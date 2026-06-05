package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/body"
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/monitors"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

func newSilencesCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "silences",
		Aliases: []string{"silence"},
		Short:   "Manage alert silences through the official Groundcover SDK",
	}
	cmd.AddCommand(silencesListCommand(cfg))
	cmd.AddCommand(silencesGetCommand(cfg))
	cmd.AddCommand(silencesCreateCommand(cfg))
	cmd.AddCommand(silencesUpdateCommand(cfg))
	cmd.AddCommand(silencesDeleteCommand(cfg))
	return cmd
}

func silencesListCommand(cfg *config.Config) *cobra.Command {
	var active bool
	var includeRecurring bool
	var limit int64
	var skip int64
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List silences",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			params := monitors.NewGetAllSilencesParams().WithContext(cmd.Context())
			if cmd.Flags().Changed("active") {
				params = params.WithActive(&active)
			}
			if cmd.Flags().Changed("include-recurring") {
				params = params.WithIncludeRecurring(&includeRecurring)
			}
			if cmd.Flags().Changed("limit") {
				params = params.WithLimit(&limit)
			}
			if cmd.Flags().Changed("skip") {
				params = params.WithSkip(&skip)
			}
			resp, err := client.Monitors.GetAllSilences(params, nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	cmd.Flags().BoolVar(&active, "active", false, "show only active silences")
	cmd.Flags().BoolVar(&includeRecurring, "include-recurring", false, "include recurring silence instances")
	cmd.Flags().Int64Var(&limit, "limit", 0, "maximum silences to return")
	cmd.Flags().Int64Var(&skip, "skip", 0, "silences to skip")
	return cmd
}

func silencesGetCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a silence",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.GetSilence(monitors.NewGetSilenceParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func silencesCreateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a silence from a JSON body",
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.CreateSilenceRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.CreateSilence(monitors.NewCreateSilenceParams().WithContext(cmd.Context()).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func silencesUpdateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a silence from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.UpdateSilenceRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.UpdateSilence(monitors.NewUpdateSilenceParams().WithContext(cmd.Context()).WithID(args[0]).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func silencesDeleteCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a silence",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.DeleteSilence(monitors.NewDeleteSilenceParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code(), "id": args[0]}, cfg.Raw)
		},
	}
}
