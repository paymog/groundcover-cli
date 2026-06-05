package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/body"
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/dashboards"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

func newDashboardsCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dashboards",
		Aliases: []string{"dashboard"},
		Short:   "Manage dashboards through the official Groundcover SDK",
	}

	cmd.AddCommand(dashboardsListCommand(cfg))
	cmd.AddCommand(dashboardsGetCommand(cfg))
	cmd.AddCommand(dashboardsCreateCommand(cfg))
	cmd.AddCommand(dashboardsUpdateCommand(cfg))
	cmd.AddCommand(dashboardsDeleteCommand(cfg))
	cmd.AddCommand(dashboardsArchiveCommand(cfg))
	cmd.AddCommand(dashboardsRestoreCommand(cfg))
	return cmd
}

func dashboardsArchiveCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "archive <id>",
		Short: "Archive a dashboard",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Dashboards.ArchiveDashboard(dashboards.NewArchiveDashboardParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code(), "id": args[0]}, cfg.Raw)
		},
	}
}

func dashboardsRestoreCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "restore <id>",
		Short: "Restore an archived dashboard",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Dashboards.RestoreDashboard(dashboards.NewRestoreDashboardParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code(), "id": args[0]}, cfg.Raw)
		},
	}
}

func dashboardsListCommand(cfg *config.Config) *cobra.Command {
	var status string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List dashboards",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			params := dashboards.NewGetDashboardsParams().WithContext(cmd.Context())
			if status != "" {
				params = params.WithStatus(&status)
			}
			resp, err := client.Dashboards.GetDashboards(params, nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	cmd.Flags().StringVar(&status, "status", "", "dashboard status filter")
	return cmd
}

func dashboardsGetCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a dashboard",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Dashboards.GetDashboard(dashboards.NewGetDashboardParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func dashboardsCreateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a dashboard from a JSON body",
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.CreateDashboardRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Dashboards.CreateDashboard(dashboards.NewCreateDashboardParams().WithContext(cmd.Context()).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func dashboardsUpdateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a dashboard from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.UpdateDashboardRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Dashboards.UpdateDashboard(dashboards.NewUpdateDashboardParams().WithContext(cmd.Context()).WithID(args[0]).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func dashboardsDeleteCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a dashboard",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Dashboards.DeleteDashboard(dashboards.NewDeleteDashboardParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func addBodyFlags(cmd *cobra.Command, input *body.Input) {
	cmd.Flags().StringVar(&input.File, "body-file", "", "JSON or YAML request body file")
	cmd.Flags().StringVar(&input.JSON, "body-json", "", "inline JSON request body")
}
