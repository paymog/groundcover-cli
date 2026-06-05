package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/body"
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/monitors"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

func newRecurringSilencesCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "recurring-silences",
		Aliases: []string{"recurring-silence"},
		Short:   "Manage recurring alert silences through the official Groundcover SDK",
	}
	cmd.AddCommand(recurringSilencesListCommand(cfg))
	cmd.AddCommand(recurringSilencesGetCommand(cfg))
	cmd.AddCommand(recurringSilencesCreateCommand(cfg))
	cmd.AddCommand(recurringSilencesUpdateCommand(cfg))
	cmd.AddCommand(recurringSilencesDeleteCommand(cfg))
	return cmd
}

func recurringSilencesListCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List recurring silences",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.GetAllRecurringSilences(monitors.NewGetAllRecurringSilencesParams().WithContext(cmd.Context()), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func recurringSilencesGetCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a recurring silence",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.GetRecurringSilence(monitors.NewGetRecurringSilenceParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func recurringSilencesCreateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a recurring silence from a JSON body",
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.CreateRecurringSilenceRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.CreateRecurringSilence(monitors.NewCreateRecurringSilenceParams().WithContext(cmd.Context()).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func recurringSilencesUpdateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a recurring silence from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.UpdateRecurringSilenceRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.UpdateRecurringSilence(monitors.NewUpdateRecurringSilenceParams().WithContext(cmd.Context()).WithID(args[0]).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func recurringSilencesDeleteCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a recurring silence",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.DeleteRecurringSilence(monitors.NewDeleteRecurringSilenceParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code(), "id": args[0]}, cfg.Raw)
		},
	}
}
