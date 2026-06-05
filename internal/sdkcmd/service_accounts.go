package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/body"
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/serviceaccounts"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

func newServiceAccountsCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service-accounts",
		Aliases: []string{"serviceaccounts", "service-account"},
		Short:   "Manage service accounts through the official Groundcover SDK",
	}
	cmd.AddCommand(serviceAccountsListCommand(cfg))
	cmd.AddCommand(serviceAccountsGetCommand(cfg))
	cmd.AddCommand(serviceAccountsCreateCommand(cfg))
	cmd.AddCommand(serviceAccountsUpdateCommand(cfg))
	cmd.AddCommand(serviceAccountsDeleteCommand(cfg))
	return cmd
}

func serviceAccountsListCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List service accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Serviceaccounts.ListServiceAccounts(serviceaccounts.NewListServiceAccountsParams().WithContext(cmd.Context()), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func serviceAccountsGetCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a service account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Serviceaccounts.GetServiceAccount(serviceaccounts.NewGetServiceAccountParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func serviceAccountsCreateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a service account from a JSON body",
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.CreateServiceAccountRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Serviceaccounts.CreateServiceAccount(serviceaccounts.NewCreateServiceAccountParams().WithContext(cmd.Context()).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func serviceAccountsUpdateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a service account from a JSON body (id lives in body)",
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.UpdateServiceAccountRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Serviceaccounts.UpdateServiceAccount(serviceaccounts.NewUpdateServiceAccountParams().WithContext(cmd.Context()).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func serviceAccountsDeleteCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a service account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Serviceaccounts.DeleteServiceAccount(serviceaccounts.NewDeleteServiceAccountParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code(), "id": args[0]}, cfg.Raw)
		},
	}
}
