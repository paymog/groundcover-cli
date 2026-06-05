package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/body"
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/monitors"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

func newMonitorsCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "monitors",
		Aliases: []string{"monitor"},
		Short:   "Manage monitors through the official Groundcover SDK",
	}
	cmd.AddCommand(monitorsListCommand(cfg))
	cmd.AddCommand(monitorsGetCommand(cfg))
	cmd.AddCommand(monitorsCreateCommand(cfg))
	cmd.AddCommand(monitorsUpdateCommand(cfg))
	cmd.AddCommand(monitorsDeleteCommand(cfg))
	return cmd
}

func monitorsListCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	var query string
	var sort string
	var limit int64
	var skip int64
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List monitors",
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &models.MonitorListRequest{}
			if input.File != "" || input.JSON != "" {
				if err := body.Decode(input, request); err != nil {
					return err
				}
			}
			if query != "" {
				request.Query = query
			}
			if sort != "" {
				request.Sort = sort
			}
			if cmd.Flags().Changed("limit") {
				request.Limit = limit
			}
			if cmd.Flags().Changed("skip") {
				request.Skip = skip
			}

			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.ListMonitors(monitors.NewListMonitorsParams().WithContext(cmd.Context()).WithBody(request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	cmd.Flags().StringVar(&query, "query", "", "GCQL monitor filter")
	cmd.Flags().StringVar(&sort, "sort", "", "sort field")
	cmd.Flags().Int64Var(&limit, "limit", 0, "maximum monitors to return")
	cmd.Flags().Int64Var(&skip, "skip", 0, "monitors to skip")
	return cmd
}

func monitorsGetCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get monitor YAML",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.GetMonitor(monitors.NewGetMonitorParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func monitorsCreateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a monitor from a YAML or JSON body",
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.CreateMonitorRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.CreateMonitor(monitors.NewCreateMonitorParams().WithContext(cmd.Context()).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func monitorsUpdateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a monitor from a YAML or JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.UpdateMonitorRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.UpdateMonitor(monitors.NewUpdateMonitorParams().WithContext(cmd.Context()).WithID(args[0]).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code(), "id": args[0]}, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func monitorsDeleteCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a monitor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Monitors.DeleteMonitor(monitors.NewDeleteMonitorParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code(), "id": args[0]}, cfg.Raw)
		},
	}
}
