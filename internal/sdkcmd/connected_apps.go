package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/body"
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/connected_apps"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

func newConnectedAppsCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connected-apps",
		Aliases: []string{"connected-app"},
		Short:   "Manage connected apps through the official Groundcover SDK",
	}
	cmd.AddCommand(connectedAppsListCommand(cfg))
	cmd.AddCommand(connectedAppsGetCommand(cfg))
	cmd.AddCommand(connectedAppsCreateCommand(cfg))
	cmd.AddCommand(connectedAppsUpdateCommand(cfg))
	cmd.AddCommand(connectedAppsDeleteCommand(cfg))
	return cmd
}

func connectedAppsListCommand(cfg *config.Config) *cobra.Command {
	var query string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List connected apps",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			params := connected_apps.NewListConnectedAppsParams().
				WithContext(cmd.Context()).
				WithBody(&models.ListConnectedAppsRequest{Query: query})
			resp, err := client.ConnectedApps.ListConnectedApps(params, nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "gcQL filter (e.g. 'type:slack-webhook')")
	return cmd
}

func connectedAppsGetCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a connected app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.ConnectedApps.GetConnectedApp(connected_apps.NewGetConnectedAppParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func connectedAppsCreateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a connected app from a JSON body",
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.CreateConnectedAppRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.ConnectedApps.CreateConnectedApp(connected_apps.NewCreateConnectedAppParams().WithContext(cmd.Context()).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func connectedAppsUpdateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a connected app from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.UpdateConnectedAppRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.ConnectedApps.UpdateConnectedApp(connected_apps.NewUpdateConnectedAppParams().WithContext(cmd.Context()).WithID(args[0]).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func connectedAppsDeleteCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a connected app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.ConnectedApps.DeleteConnectedApp(connected_apps.NewDeleteConnectedAppParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code(), "id": args[0]}, cfg.Raw)
		},
	}
}
