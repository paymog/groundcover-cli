package sdkcmd

import (
	"fmt"
	"os"

	"github.com/paymog/groundcover-cli/internal/body"
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/workflows"
	"github.com/spf13/cobra"
)

func newWorkflowsCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "workflows",
		Aliases: []string{"workflow"},
		Short:   "Manage workflows through the official Groundcover SDK",
	}
	cmd.AddCommand(workflowsListCommand(cfg))
	cmd.AddCommand(workflowsCreateCommand(cfg))
	cmd.AddCommand(workflowsDeleteCommand(cfg))
	return cmd
}

func workflowsListCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List workflows",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Workflows.ListWorkflows(workflows.NewListWorkflowsParams().WithContext(cmd.Context()), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
}

func workflowsCreateCommand(cfg *config.Config) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a workflow from a raw body file (sent as text/plain)",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := readRawBody(input)
			if err != nil {
				return err
			}
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Workflows.CreateWorkflow(workflows.NewCreateWorkflowParams().WithContext(cmd.Context()).WithBody(raw), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), resp.Payload, cfg.Raw)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func workflowsDeleteCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a workflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.Workflows.DeleteWorkflow(workflows.NewDeleteWorkflowParams().WithContext(cmd.Context()).WithID(args[0]), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code(), "id": args[0]}, cfg.Raw)
		},
	}
}

func readRawBody(input body.Input) (string, error) {
	if input.File != "" && input.JSON != "" {
		return "", fmt.Errorf("use only one of --body-file or --body-json")
	}
	if input.JSON != "" {
		return input.JSON, nil
	}
	if input.File == "" {
		return "", fmt.Errorf("missing request body: pass --body-file or --body-json")
	}
	data, err := os.ReadFile(input.File)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
