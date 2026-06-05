package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/body"
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/aggregations_metrics"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/logs_pipeline"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics_pipeline"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/traces_pipeline"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

// Pipeline configs share the same SDK shape: get/create/update/delete on a singleton config,
// where the non-GET verbs all consume the same CreateOrUpdateXRequest body.

// pipelineMutationCommand wires create/update/delete for a singleton pipeline config.
// invoke is responsible for decoding the body and calling the right SDK method.
func pipelineMutationCommand[T any](use, short string, invoke func(*cobra.Command, *T, body.Input) (any, error)) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			var request T
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			payload, err := invoke(cmd, &request, input)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), payload, false)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func newLogsPipelineCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{Use: "logs-pipeline", Short: "Manage the logs pipeline config through the official Groundcover SDK"}
	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get the logs pipeline config",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			ok, noContent, err := client.LogsPipeline.GetLogsPipelineConfig(logs_pipeline.NewGetLogsPipelineConfigParams().WithContext(cmd.Context()), nil)
			if err != nil {
				return err
			}
			if ok != nil {
				return output.Print(cmd.OutOrStdout(), ok.Payload, cfg.Raw)
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": noContent.Code()}, cfg.Raw)
		},
	})
	cmd.AddCommand(pipelineMutationCommand("create", "Create the logs pipeline config from a JSON body",
		func(c *cobra.Command, req *models.CreateOrUpdateLogsPipelineConfigRequest, _ body.Input) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.LogsPipeline.CreateLogsPipelineConfig(logs_pipeline.NewCreateLogsPipelineConfigParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(pipelineMutationCommand("update", "Update the logs pipeline config from a JSON body",
		func(c *cobra.Command, req *models.CreateOrUpdateLogsPipelineConfigRequest, _ body.Input) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.LogsPipeline.UpdateLogsPipelineConfig(logs_pipeline.NewUpdateLogsPipelineConfigParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete the logs pipeline config",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.LogsPipeline.DeleteLogsPipelineConfig(logs_pipeline.NewDeleteLogsPipelineConfigParams().WithContext(cmd.Context()), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code()}, cfg.Raw)
		},
	})
	return cmd
}

func newMetricsPipelineCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{Use: "metrics-pipeline", Short: "Manage the metrics pipeline config through the official Groundcover SDK"}
	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get the metrics pipeline config",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			ok, noContent, err := client.MetricsPipeline.GetMetricsPipelineConfig(metrics_pipeline.NewGetMetricsPipelineConfigParams().WithContext(cmd.Context()), nil)
			if err != nil {
				return err
			}
			if ok != nil {
				return output.Print(cmd.OutOrStdout(), ok.Payload, cfg.Raw)
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": noContent.Code()}, cfg.Raw)
		},
	})
	cmd.AddCommand(pipelineMutationCommand("create", "Create the metrics pipeline config from a JSON body",
		func(c *cobra.Command, req *models.CreateOrUpdateMetricsPipelineConfigRequest, _ body.Input) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.MetricsPipeline.CreateMetricsPipelineConfig(metrics_pipeline.NewCreateMetricsPipelineConfigParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(pipelineMutationCommand("update", "Update the metrics pipeline config from a JSON body",
		func(c *cobra.Command, req *models.CreateOrUpdateMetricsPipelineConfigRequest, _ body.Input) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.MetricsPipeline.UpdateMetricsPipelineConfig(metrics_pipeline.NewUpdateMetricsPipelineConfigParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete the metrics pipeline config",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.MetricsPipeline.DeleteMetricsPipelineConfig(metrics_pipeline.NewDeleteMetricsPipelineConfigParams().WithContext(cmd.Context()), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code()}, cfg.Raw)
		},
	})
	return cmd
}

func newTracesPipelineCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{Use: "traces-pipeline", Short: "Manage the traces pipeline config through the official Groundcover SDK"}
	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get the traces pipeline config",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			ok, noContent, err := client.TracesPipeline.GetTracesPipelineConfig(traces_pipeline.NewGetTracesPipelineConfigParams().WithContext(cmd.Context()), nil)
			if err != nil {
				return err
			}
			if ok != nil {
				return output.Print(cmd.OutOrStdout(), ok.Payload, cfg.Raw)
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": noContent.Code()}, cfg.Raw)
		},
	})
	cmd.AddCommand(pipelineMutationCommand("create", "Create the traces pipeline config from a JSON body",
		func(c *cobra.Command, req *models.CreateOrUpdateTracesPipelineConfigRequest, _ body.Input) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.TracesPipeline.CreateTracesPipelineConfig(traces_pipeline.NewCreateTracesPipelineConfigParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(pipelineMutationCommand("update", "Update the traces pipeline config from a JSON body",
		func(c *cobra.Command, req *models.CreateOrUpdateTracesPipelineConfigRequest, _ body.Input) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.TracesPipeline.UpdateTracesPipelineConfig(traces_pipeline.NewUpdateTracesPipelineConfigParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete the traces pipeline config",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.TracesPipeline.DeleteTracesPipelineConfig(traces_pipeline.NewDeleteTracesPipelineConfigParams().WithContext(cmd.Context()), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code()}, cfg.Raw)
		},
	})
	return cmd
}

func newMetricsAggregatorCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{Use: "metrics-aggregator", Short: "Manage the metrics aggregator config through the official Groundcover SDK"}
	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get the metrics aggregator config",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			ok, noContent, err := client.AggregationsMetrics.GetMetricsAggregatorConfig(aggregations_metrics.NewGetMetricsAggregatorConfigParams().WithContext(cmd.Context()), nil)
			if err != nil {
				return err
			}
			if ok != nil {
				return output.Print(cmd.OutOrStdout(), ok.Payload, cfg.Raw)
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": noContent.Code()}, cfg.Raw)
		},
	})
	cmd.AddCommand(pipelineMutationCommand("create", "Create the metrics aggregator config from a JSON body",
		func(c *cobra.Command, req *models.CreateOrUpdateMetricsAggregatorConfigRequest, _ body.Input) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.AggregationsMetrics.CreateMetricsAggregatorConfig(aggregations_metrics.NewCreateMetricsAggregatorConfigParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(pipelineMutationCommand("update", "Update the metrics aggregator config from a JSON body",
		func(c *cobra.Command, req *models.CreateOrUpdateMetricsAggregatorConfigRequest, _ body.Input) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.AggregationsMetrics.UpdateMetricsAggregatorConfig(aggregations_metrics.NewUpdateMetricsAggregatorConfigParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete the metrics aggregator config",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(*cfg)
			if err != nil {
				return err
			}
			resp, err := client.AggregationsMetrics.DeleteMetricsAggregatorConfig(aggregations_metrics.NewDeleteMetricsAggregatorConfigParams().WithContext(cmd.Context()), nil)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), map[string]any{"status": resp.Code()}, cfg.Raw)
		},
	})
	return cmd
}
