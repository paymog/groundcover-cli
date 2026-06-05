package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/body"
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/integrations"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

func newIntegrationsCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "integrations",
		Aliases: []string{"integration"},
		Short:   "Manage data integration configs through the official Groundcover SDK",
	}
	cmd.AddCommand(integrationsListCommand(cfg))
	cmd.AddCommand(integrationsListByTypeCommand(cfg))
	cmd.AddCommand(integrationsGetCommand(cfg))
	cmd.AddCommand(integrationsDescribeCommand(cfg))
	cmd.AddCommand(integrationsCreateCommand(cfg))
	cmd.AddCommand(integrationsUpdateCommand(cfg))
	cmd.AddCommand(integrationsDeleteCommand(cfg))
	return cmd
}

func integrationsListCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all data integration configs",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Integrations.GetDataIntegrationConfigs(integrations.NewGetDataIntegrationConfigsParams().WithContext(cmd.Context()), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func integrationsListByTypeCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list-by-type <type>",
		Short: "List data integration configs of a given type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Integrations.GetDataIntegrationConfigsByType(integrations.NewGetDataIntegrationConfigsByTypeParams().WithContext(cmd.Context()).WithType(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func integrationsGetCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <type> <id>",
		Short: "Get a data integration config",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Integrations.GetDataIntegrationConfig(integrations.NewGetDataIntegrationConfigParams().WithContext(cmd.Context()).WithType(args[0]).WithID(args[1]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func integrationsDescribeCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "describe <type>",
		Short: "Describe a data integration type (schema)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Integrations.DescribeDataIntegration(integrations.NewDescribeDataIntegrationParams().WithContext(cmd.Context()).WithType(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func integrationsCreateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "create <type>",
		Short: "Create a data integration config of the given type from a JSON body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.CreateDataIntegrationConfigRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Integrations.CreateDataIntegrationConfig(integrations.NewCreateDataIntegrationConfigParams().WithContext(cmd.Context()).WithType(args[0]).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func integrationsUpdateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "update <type> <id>",
		Short: "Update a data integration config from a JSON body",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.CreateDataIntegrationConfigRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Integrations.UpdateDataIntegrationConfig(integrations.NewUpdateDataIntegrationConfigParams().WithContext(cmd.Context()).WithType(args[0]).WithID(args[1]).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func integrationsDeleteCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <type> <id>",
		Short: "Delete a data integration config",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Integrations.DeleteDataIntegrationConfig(integrations.NewDeleteDataIntegrationConfigParams().WithContext(cmd.Context()).WithType(args[0]).WithID(args[1]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code(), "type": args[0], "id": args[1]}, cfg.Raw)
		},
	}
}
