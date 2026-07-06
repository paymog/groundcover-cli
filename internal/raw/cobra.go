package raw

import (
	"fmt"
	"strings"

	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:                "raw <command...>",
		Short:              "Run HAR-derived webapp endpoints not covered by the SDK",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		SilenceUsage:       true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
				return cmd.Help()
			}
			if args[0] == "list" {
				filter := strings.Join(args[1:], " ")
				printCommands(cmd, filter)
				return nil
			}

			tokens, flagArgs := splitArgs(args)
			command, ok := Find(tokens)
			if !ok {
				return fmt.Errorf("unknown raw command %q; run `groundcover raw list`", strings.Join(tokens, " "))
			}

			localCfg := *cfg
			localCfg.ApplyEnv()
			opts, err := parseOptions(flagArgs, command, &localCfg)
			if err != nil {
				return err
			}
			return Run(command, localCfg, opts, cmd.OutOrStdout())
		},
	}
	return cmd
}

func printCommands(cmd *cobra.Command, filter string) {
	for _, command := range allCommands() {
		name := command.Key()
		if filter != "" && !strings.Contains(name, filter) {
			continue
		}
		params := make([]string, 0, len(command.PathParams))
		for _, param := range command.PathParams {
			params = append(params, fmt.Sprintf("--%s <%s>", kebab(param), param))
		}
		if len(params) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s\t%s\n", name, strings.Join(params, " "), command.Description)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", name, command.Description)
		}
	}
}
