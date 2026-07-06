package main

import (
	"reflect"
	"testing"
)

func TestNormalizedPathGrafanaDashboardUID(t *testing.T) {
	path, params := normalizedPath("/grafana/api/dashboards/uid/streamling-pipeline-slo")
	if path != "/grafana/api/dashboards/uid/:dashboardUid" {
		t.Fatalf("unexpected path %s", path)
	}
	if !reflect.DeepEqual(params, []string{"dashboardUid"}) {
		t.Fatalf("unexpected params %#v", params)
	}

	name := commandName("GET", path)
	if !reflect.DeepEqual(name, []string{"grafana", "dashboards", "get"}) {
		t.Fatalf("unexpected name %#v", name)
	}
}

func TestNormalizedPathGrafanaFolderUID(t *testing.T) {
	path, params := normalizedPath("/grafana/api/folders/bend1nm1f0ruod")
	if path != "/grafana/api/folders/:folderUid" {
		t.Fatalf("unexpected path %s", path)
	}
	if !reflect.DeepEqual(params, []string{"folderUid"}) {
		t.Fatalf("unexpected params %#v", params)
	}

	name := commandName("GET", path)
	if !reflect.DeepEqual(name, []string{"grafana", "folders", "get"}) {
		t.Fatalf("unexpected name %#v", name)
	}
}

func TestNormalizedPathGrafanaDatasourceLabelValues(t *testing.T) {
	path, params := normalizedPath("/grafana/api/datasources/uid/aelovgen78268b/resources/api/v1/label/project_id/values")
	if path != "/grafana/api/datasources/uid/:datasourceUid/resources/api/v1/label/:label/values" {
		t.Fatalf("unexpected path %s", path)
	}
	if !reflect.DeepEqual(params, []string{"datasourceUid", "label"}) {
		t.Fatalf("unexpected params %#v", params)
	}

	name := commandName("GET", path)
	want := []string{"grafana", "datasources", "resources", "api", "label", "values"}
	if !reflect.DeepEqual(name, want) {
		t.Fatalf("unexpected name %#v", name)
	}
}
