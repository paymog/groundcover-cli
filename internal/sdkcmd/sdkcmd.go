package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/spf13/cobra"
)

func AddCommands(root *cobra.Command, cfg *config.Config) {
	root.AddCommand(newDashboardsCommand(cfg))
	root.AddCommand(newMonitorsCommand(cfg))
	root.AddCommand(newSilencesCommand(cfg))
	root.AddCommand(newRecurringSilencesCommand(cfg))
	root.AddCommand(newConnectedAppsCommand(cfg))
	root.AddCommand(newNotificationRoutesCommand(cfg))
	root.AddCommand(newAPIKeysCommand(cfg))
	root.AddCommand(newIngestionKeysCommand(cfg))
	root.AddCommand(newServiceAccountsCommand(cfg))
	root.AddCommand(newPoliciesCommand(cfg))
	root.AddCommand(newSyntheticsCommand(cfg))
	root.AddCommand(newSecretsCommand(cfg))
	root.AddCommand(newWorkflowsCommand(cfg))
	root.AddCommand(newIntegrationsCommand(cfg))
	root.AddCommand(newLogsPipelineCommand(cfg))
	root.AddCommand(newMetricsPipelineCommand(cfg))
	root.AddCommand(newTracesPipelineCommand(cfg))
	root.AddCommand(newMetricsAggregatorCommand(cfg))
	root.AddCommand(newLogsCommand(cfg))
	root.AddCommand(newTracesCommand(cfg))
	root.AddCommand(newMetricsCommand(cfg))
	root.AddCommand(newSearchCommand(cfg))
	root.AddCommand(newK8sCommand(cfg))
}
