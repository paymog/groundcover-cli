package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/body"
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/k8s"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/logs"
	metricsclient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/metrics"
	searchclient "github.com/groundcover-com/groundcover-sdk-go/pkg/client/search"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/traces"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

// bodyQueryCommand wires a single body-only POST endpoint into a leaf command.
func bodyQueryCommand[Req any](use, short string, invoke func(cmd *cobra.Command, req *Req) (any, error)) *cobra.Command {
	input := body.Input{}
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			var request Req
			if err := body.Decode(input, &request); err != nil {
				return err
			}
			payload, err := invoke(cmd, &request)
			if err != nil {
				return err
			}
			return output.Print(cmd.OutOrStdout(), payload, false)
		},
	}
	addBodyFlags(cmd, &input)
	return cmd
}

func newLogsCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{Use: "logs", Short: "Query logs through the official Groundcover SDK"}
	cmd.AddCommand(bodyQueryCommand("search", "Search logs (POST /api/logs/v2/search)",
		func(c *cobra.Command, req *models.LogsSearchRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.Logs.SearchLogs(logs.NewSearchLogsParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	return cmd
}

func newTracesCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{Use: "traces", Short: "Query traces through the official Groundcover SDK"}
	cmd.AddCommand(bodyQueryCommand("search", "Search traces (POST /api/traces/v2/search)",
		func(c *cobra.Command, req *models.TracesSearchRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.Traces.SearchTraces(traces.NewSearchTracesParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	return cmd
}

func newMetricsCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{Use: "metrics", Short: "Query metrics through the official Groundcover SDK"}
	cmd.AddCommand(bodyQueryCommand("query", "Run a metrics query (POST /api/metrics/query)",
		func(c *cobra.Command, req *models.QueryRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.Metrics.MetricsQuery(metricsclient.NewMetricsQueryParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(bodyQueryCommand("names", "List metric names (POST /api/metrics/names)",
		func(c *cobra.Command, req *models.MetricsNamesRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.Metrics.GetMetricNames(metricsclient.NewGetMetricNamesParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(bodyQueryCommand("keys", "List metric keys (POST /api/metrics/keys)",
		func(c *cobra.Command, req *models.MetricsKeysRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.Metrics.GetMetricKeys(metricsclient.NewGetMetricKeysParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(bodyQueryCommand("values", "List metric values (POST /api/metrics/values)",
		func(c *cobra.Command, req *models.MetricsValuesRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.Metrics.GetMetricValues(metricsclient.NewGetMetricValuesParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	return cmd
}

func newSearchCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{Use: "search", Short: "Run search-API queries through the official Groundcover SDK"}
	cmd.AddCommand(bodyQueryCommand("discovery", "Run a discovery query (POST /api/search/discovery)",
		func(c *cobra.Command, req *models.DiscoveryRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.Search.GetDiscovery(searchclient.NewGetDiscoveryParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(bodyQueryCommand("keys", "Search keys (POST /api/search/keys)",
		func(c *cobra.Command, req *models.KeysRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.Search.GetKeys(searchclient.NewGetKeysParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(bodyQueryCommand("values", "Search values (POST /api/search/values)",
		func(c *cobra.Command, req *models.ValuesRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.Search.GetValues(searchclient.NewGetValuesParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	return cmd
}

func newK8sCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{Use: "k8s", Short: "Query Kubernetes inventory and events through the official Groundcover SDK"}
	cmd.AddCommand(bodyQueryCommand("clusters", "List clusters (POST /api/k8s/v3/clusters/list)",
		func(c *cobra.Command, req *models.ClustersListRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.K8s.ClustersList(k8s.NewClustersListParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(bodyQueryCommand("workloads", "List workloads (POST /api/k8s/v3/workloads/list)",
		func(c *cobra.Command, req *models.WorkloadsListRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.K8s.WorkloadsList(k8s.NewWorkloadsListParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(bodyQueryCommand("events-search", "Search Kubernetes events (POST /api/k8s/v2/events/search)",
		func(c *cobra.Command, req *models.EventsSearchRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.K8s.EventsSearch(k8s.NewEventsSearchParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	cmd.AddCommand(bodyQueryCommand("events-over-time", "Aggregate Kubernetes events over time (POST /api/k8s/v2/events-over-time)",
		func(c *cobra.Command, req *models.GetEventsOverTimeRequest) (any, error) {
			client, err := newClient(*cfg)
			if err != nil {
				return nil, err
			}
			resp, err := client.K8s.GetEventsOverTime(k8s.NewGetEventsOverTimeParams().WithContext(c.Context()).WithBody(req), nil)
			if err != nil {
				return nil, err
			}
			return resp.Payload, nil
		}))
	return cmd
}
