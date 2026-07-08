package raw

import (
	"strings"
	"testing"
)

func TestGrafanaSetupGuideContent(t *testing.T) {
	guide := grafanaSetupGuide()
	required := []string{
		"https://groundcover.com/install.sh",
		"auth login",
		"generate-service-account-token",
		"~/.groundcover/bin",
		"GROUNDCOVER_GRAFANA_SERVICE_ACCOUNT_TOKEN",
		"PATH",
	}
	for _, want := range required {
		if !strings.Contains(guide, want) {
			t.Errorf("grafanaSetupGuide() missing substring %q", want)
		}
	}
}
