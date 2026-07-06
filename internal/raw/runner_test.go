package raw

import (
	"testing"

	"github.com/paymog/groundcover-cli/internal/config"
)

func TestSetDeep(t *testing.T) {
	target := map[string]any{}
	setDeep(target, "a.b.c", 1)

	a, ok := target["a"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map at a")
	}
	b, ok := a["b"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map at a.b")
	}
	if b["c"] != 1 {
		t.Fatalf("expected c=1, got %#v", b["c"])
	}
}

func TestFind(t *testing.T) {
	command, ok := Find([]string{"backend", "settings"})
	if !ok {
		t.Fatalf("expected backend settings command")
	}
	if command.Path != "/api/backend/settings" {
		t.Fatalf("unexpected path %s", command.Path)
	}
}

func TestFindGrafanaCommand(t *testing.T) {
	command, ok := Find([]string{"grafana", "dashboards", "get"})
	if !ok {
		t.Fatalf("expected grafana dashboards get command")
	}
	if command.Path != "/grafana/api/dashboards/uid/:dashboardUid" {
		t.Fatalf("unexpected path %s", command.Path)
	}
	if len(command.PathParams) != 1 || command.PathParams[0] != "dashboardUid" {
		t.Fatalf("unexpected path params %#v", command.PathParams)
	}
	if !command.WebApp {
		t.Fatalf("expected grafana command to target webapp host")
	}
}

func TestBuildURLUsesWebAppDefault(t *testing.T) {
	command, ok := Find([]string{"grafana", "dashboards", "get"})
	if !ok {
		t.Fatalf("expected grafana dashboards get command")
	}
	requestURL, err := buildURL(command, config.Config{BaseURL: config.DefaultBaseURL}, Options{
		PathValues: map[string]string{"dashboardUid": "streamling-pipeline-slo"},
		Query:      []string{"orgId=1"},
	})
	if err != nil {
		t.Fatalf("buildURL failed: %v", err)
	}
	if got, want := requestURL.String(), "https://app.groundcover.com/grafana/api/dashboards/uid/streamling-pipeline-slo?orgId=1"; got != want {
		t.Fatalf("unexpected URL\n got: %s\nwant: %s", got, want)
	}
}

func TestBuildURLHonorsCustomBaseURL(t *testing.T) {
	command, ok := Find([]string{"grafana", "folders", "get"})
	if !ok {
		t.Fatalf("expected grafana folders get command")
	}
	requestURL, err := buildURL(command, config.Config{BaseURL: "https://groundcover.example"}, Options{
		PathValues: map[string]string{"folderUid": "bend1nm1f0ruod"},
	})
	if err != nil {
		t.Fatalf("buildURL failed: %v", err)
	}
	if got, want := requestURL.String(), "https://groundcover.example/grafana/api/folders/bend1nm1f0ruod"; got != want {
		t.Fatalf("unexpected URL\n got: %s\nwant: %s", got, want)
	}
}
