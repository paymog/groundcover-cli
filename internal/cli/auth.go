package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/credstore"
	"github.com/paymog/groundcover-cli/internal/sdkcmd"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// skipResolveAnnotation marks commands that must run before credentials are
// resolved (e.g. `auth login`, which has no stored credentials yet). The root
// PersistentPreRunE skips config.Resolve for any command carrying it.
const skipResolveAnnotation = "skipAuthResolve"

func skipResolve() map[string]string {
	return map[string]string{skipResolveAnnotation: "true"}
}

func newAuthCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage stored credential profiles",
		Long: "Manage named credential profiles. API keys are stored in the OS keyring;\n" +
			"profile metadata (backend ID, base URL, tenant UUID) lives in a config file.\n\n" +
			"For CI or one-off use, set GROUNDCOVER_API_KEY instead — it takes precedence\n" +
			"over stored profiles.",
		Annotations: skipResolve(),
		RunE:        func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
	}
	cmd.AddCommand(
		authLoginCommand(cfg),
		authListCommand(),
		authDefaultCommand(),
		authLogoutCommand(),
		authTokenCommand(cfg),
		authStatusCommand(cfg),
	)
	return cmd
}

func authLoginCommand(cfg *config.Config) *cobra.Command {
	var (
		name      string
		key       string
		backendID string
		baseURL   string
		tenant    string
	)
	cmd := &cobra.Command{
		Use:         "login [name]",
		Short:       "Add or update a credential profile",
		Args:        cobra.MaximumNArgs(1),
		Annotations: skipResolve(),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				name = args[0]
			}
			if name == "" {
				name = "default"
			}
			if !credstore.Available() {
				return fmt.Errorf("no usable OS keyring found; set GROUNDCOVER_API_KEY instead of using profiles")
			}

			if key == "" {
				var err error
				key, err = promptSecret(cmd, "Enter your Groundcover API key: ")
				if err != nil {
					return err
				}
			}
			key = strings.TrimSpace(key)
			if key == "" {
				return fmt.Errorf("no API key provided")
			}

			if backendID == "" {
				backendID = config.DefaultBackendID
			}

			// Validate the credentials before persisting them.
			probe := config.Config{APIKey: key, BackendID: backendID, BaseURL: baseURL, Timeout: cfg.Timeout}
			if err := sdkcmd.ValidateAuth(cmd.Context(), probe); err != nil {
				return err
			}

			store, err := credstore.Load()
			if err != nil {
				return err
			}
			first := len(store.Profiles) == 0
			if err := store.Add(name, credstore.Profile{
				BackendID:  backendID,
				BaseURL:    baseURL,
				TenantUUID: tenant,
			}, key); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Saved profile %q (backend %s)\n", name, backendID)
			if first {
				fmt.Fprintln(cmd.OutOrStdout(), "  Set as default profile")
			}
			if os.Getenv("GROUNDCOVER_API_KEY") != "" || os.Getenv("GC_API_KEY") != "" {
				fmt.Fprintln(cmd.ErrOrStderr(), "Warning: GROUNDCOVER_API_KEY is set and takes precedence over stored profiles.")
			}
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVar(&name, "name", "", "Profile name (defaults to the positional arg, then \"default\")")
	f.StringVarP(&key, "key", "k", "", "API key (prompted if not provided)")
	f.StringVar(&backendID, "backend-id", "", "Backend ID for this profile (default \""+config.DefaultBackendID+"\")")
	f.StringVar(&baseURL, "base-url", "", "Base URL for this profile (optional)")
	f.StringVar(&tenant, "tenant-uuid", "", "Tenant UUID for this profile (optional)")
	return cmd
}

func authListCommand() *cobra.Command {
	return &cobra.Command{
		Use:         "list",
		Short:       "List configured profiles",
		Annotations: skipResolve(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			store, err := credstore.Load()
			if err != nil {
				return err
			}
			names := store.Names()
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No profiles configured. Run `groundcover auth login`.")
				return nil
			}
			tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(tw, "\tPROFILE\tBACKEND\tBASE URL")
			for _, n := range names {
				marker := " "
				if n == store.Default {
					marker = "*"
				}
				p := store.Profiles[n]
				base := p.BaseURL
				if base == "" {
					base = config.DefaultBaseURL
				}
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", marker, n, p.BackendID, base)
			}
			return tw.Flush()
		},
	}
}

func authDefaultCommand() *cobra.Command {
	return &cobra.Command{
		Use:         "default <name>",
		Short:       "Set the default profile",
		Args:        cobra.ExactArgs(1),
		Annotations: skipResolve(),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := credstore.Load()
			if err != nil {
				return err
			}
			if err := store.SetDefault(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Default profile set to %q\n", args[0])
			return nil
		},
	}
}

func authLogoutCommand() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:         "logout <name>",
		Short:       "Remove a credential profile",
		Args:        cobra.ExactArgs(1),
		Annotations: skipResolve(),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			store, err := credstore.Load()
			if err != nil {
				return err
			}
			if _, ok := store.Profiles[name]; !ok {
				return fmt.Errorf("profile %q not found", name)
			}
			if !force {
				ok, err := confirm(cmd, fmt.Sprintf("Remove profile %q? [y/N] ", name))
				if err != nil {
					return err
				}
				if !ok {
					fmt.Fprintln(cmd.OutOrStdout(), "Aborted")
					return nil
				}
			}
			if err := store.Remove(name); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed profile %q\n", name)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
	return cmd
}

func authTokenCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "token",
		Short: "Print the resolved API key",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := cfg.RequireAPIKey(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), cfg.APIKey)
			return nil
		},
	}
}

func authStatusCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the resolved credential source",
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cmd.OutOrStdout()
			switch {
			case os.Getenv("GROUNDCOVER_API_KEY") != "" || os.Getenv("GC_API_KEY") != "":
				fmt.Fprintln(out, "Source: GROUNDCOVER_API_KEY environment variable")
			case cfg.Profile != "":
				fmt.Fprintf(out, "Source: profile %q (--profile)\n", cfg.Profile)
			default:
				store, err := credstore.Load()
				if err != nil {
					return err
				}
				if store.Default != "" {
					fmt.Fprintf(out, "Source: default profile %q\n", store.Default)
				} else {
					fmt.Fprintln(out, "Source: none (no profile and no GROUNDCOVER_API_KEY)")
				}
			}
			fmt.Fprintf(out, "Backend ID: %s\n", cfg.BackendID)
			fmt.Fprintf(out, "Base URL: %s\n", cfg.NormalizedBaseURL())
			if err := cfg.RequireAPIKey(); err != nil {
				fmt.Fprintf(out, "API key: (not configured)\n")
			} else {
				fmt.Fprintf(out, "API key: configured\n")
			}
			fmt.Fprintf(out, "Keyring available: %t\n", credstore.Available())
			return nil
		},
	}
}

func promptSecret(cmd *cobra.Command, prompt string) (string, error) {
	fmt.Fprint(cmd.ErrOrStderr(), prompt)
	if term.IsTerminal(int(os.Stdin.Fd())) {
		b, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(cmd.ErrOrStderr())
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	// Non-interactive: read a line from stdin.
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && line == "" {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func confirm(cmd *cobra.Command, prompt string) (bool, error) {
	fmt.Fprint(cmd.ErrOrStderr(), prompt)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && line == "" {
		return false, err
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes", nil
}
