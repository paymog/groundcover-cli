package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/credstore"
	"github.com/paymog/groundcover-cli/internal/raw"
	"github.com/paymog/groundcover-cli/internal/sdkcmd"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cfg := config.FromEnv()

	root := &cobra.Command{
		Use:           "groundcover",
		Short:         "Groundcover CLI with SDK-backed writes and raw webapp endpoint passthrough",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Commands that manage credentials (e.g. `auth login`) run before
			// resolution; everything else resolves the active profile/env.
			if cmd.Annotations[skipResolveAnnotation] == "true" {
				cfg.ApplyEnv()
				return nil
			}
			store, err := credstore.Load()
			if err != nil {
				return err
			}
			return cfg.Resolve(store)
		},
	}

	flags := root.PersistentFlags()
	flags.StringVar(&cfg.APIKey, "api-key", cfg.APIKey, "Groundcover API key. Env: GROUNDCOVER_API_KEY or GC_API_KEY")
	flags.StringVar(&cfg.BackendID, "backend-id", cfg.BackendID, "Groundcover backend ID. Env: GROUNDCOVER_BACKEND_ID or GC_BACKEND_ID")
	flags.StringVar(&cfg.TenantUUID, "tenant-uuid", cfg.TenantUUID, "Groundcover tenant UUID (raw endpoints only). Env: GROUNDCOVER_TENANT_UUID or GC_TENANT_UUID")
	flags.StringVar(&cfg.BaseURL, "base-url", cfg.BaseURL, "Groundcover API base URL. Env: GROUNDCOVER_BASE_URL or GC_BASE_URL")
	flags.StringVar(&cfg.Profile, "profile", os.Getenv("GROUNDCOVER_PROFILE"), "Named credential profile to use. Env: GROUNDCOVER_PROFILE")
	flags.BoolVar(&cfg.Raw, "raw", false, "print response bytes without JSON formatting where supported")
	flags.DurationVar(&cfg.Timeout, "timeout", 30*time.Second, "request timeout")

	root.AddCommand(newAuthCommand(&cfg))
	root.AddCommand(raw.NewCommand(&cfg))
	sdkcmd.AddCommands(root, &cfg)

	root.SetErr(os.Stderr)
	root.SetOut(os.Stdout)
	root.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return fmt.Errorf("%w\n\nRun `%s --help` for usage", err, cmd.CommandPath())
	})

	return root
}
