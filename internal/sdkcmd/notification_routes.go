package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/body"
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/notification_routes"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

func newNotificationRoutesCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "notification-routes",
		Aliases: []string{"notification-route"},
		Short:   "Manage notification routes through the official Groundcover SDK",
	}
	cmd.AddCommand(notificationRoutesListCommand(cfg))
	cmd.AddCommand(notificationRoutesGetCommand(cfg))
	cmd.AddCommand(notificationRoutesCreateCommand(cfg))
	cmd.AddCommand(notificationRoutesUpdateCommand(cfg))
	cmd.AddCommand(notificationRoutesDeleteCommand(cfg))
	return cmd
}

func notificationRoutesListCommand(cfg *config.Config) *cobra.Command {
	var query string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List notification routes",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			params := notification_routes.NewListNotificationRoutesParams().
				WithContext(cmd.Context()).
				WithBody(&models.ListNotificationRoutesRequest{Query: query})
			resp, err := client.NotificationRoutes.ListNotificationRoutes(params, nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "freetext filter on route name")
	return cmd
}

func notificationRoutesGetCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a notification route",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.NotificationRoutes.GetNotificationRoute(notification_routes.NewGetNotificationRouteParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func notificationRoutesCreateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a notification route from a JSON body",
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.CreateNotificationRouteRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.NotificationRoutes.CreateNotificationRoute(notification_routes.NewCreateNotificationRouteParams().WithContext(cmd.Context()).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func notificationRoutesUpdateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a notification route from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.UpdateNotificationRouteRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.NotificationRoutes.UpdateNotificationRoute(notification_routes.NewUpdateNotificationRouteParams().WithContext(cmd.Context()).WithID(args[0]).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func notificationRoutesDeleteCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a notification route",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.NotificationRoutes.DeleteNotificationRoute(notification_routes.NewDeleteNotificationRouteParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code(), "id": args[0]}, cfg.Raw)
		},
	}
}
