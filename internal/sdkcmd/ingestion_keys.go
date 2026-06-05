package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/body"
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/ingestionkeys"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

func newIngestionKeysCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ingestion-keys",
		Aliases: []string{"ingestionkeys", "ingestion-key"},
		Short:   "Manage ingestion keys through the official Groundcover SDK",
	}
	cmd.AddCommand(ingestionKeysListCommand(cfg))
	cmd.AddCommand(ingestionKeysCreateCommand(cfg))
	cmd.AddCommand(ingestionKeysDeleteCommand(cfg))
	return cmd
}

func ingestionKeysListCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List ingestion keys (POST /api/rbac/ingestion-keys/list)",
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.ListIngestionKeysRequest
			if err := decodeOptionalBody(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Ingestionkeys.ListIngestionKeys(ingestionkeys.NewListIngestionKeysParams().WithContext(cmd.Context()).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func ingestionKeysCreateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an ingestion key from a JSON body",
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.CreateIngestionKeyRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Ingestionkeys.CreateIngestionKey(ingestionkeys.NewCreateIngestionKeyParams().WithContext(cmd.Context()).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func ingestionKeysDeleteCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an ingestion key (body identifies which)",
		RunE: func(cmd *cobra.Command, args []string) error {
			var request models.DeleteIngestionKeyRequest
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Ingestionkeys.DeleteIngestionKey(ingestionkeys.NewDeleteIngestionKeyParams().WithContext(cmd.Context()).WithBody(&request), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code()}, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}
