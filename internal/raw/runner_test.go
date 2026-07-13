package raw

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

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

func TestStorageManagementCommands(t *testing.T) {
	tests := []struct {
		name   string
		method string
	}{
		{name: "get", method: http.MethodGet},
		{name: "update", method: http.MethodPut},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, ok := Find([]string{"storage-management", tt.name})
			if !ok {
				t.Fatalf("expected storage-management %s command", tt.name)
			}
			if command.Method != tt.method {
				t.Fatalf("unexpected method %s", command.Method)
			}
			if command.Path != "/api/storage-management/:dataType" {
				t.Fatalf("unexpected path %s", command.Path)
			}
			if len(command.PathParams) != 1 || command.PathParams[0] != "dataType" {
				t.Fatalf("unexpected path params %#v", command.PathParams)
			}

			requestURL, err := buildURL(command, config.Config{BaseURL: config.DefaultBaseURL}, Options{
				PathValues: map[string]string{"dataType": "monitor_instance"},
			})
			if err != nil {
				t.Fatalf("buildURL failed: %v", err)
			}
			if got, want := requestURL.String(), "https://api.groundcover.com/api/storage-management/monitor_instance"; got != want {
				t.Fatalf("unexpected URL\n got: %s\nwant: %s", got, want)
			}
		})
	}
}

func TestRunStorageManagementUpdate(t *testing.T) {
	var captured *http.Request
	var capturedBody []byte
	var readErr error
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Clone(r.Context())
		capturedBody, readErr = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	command, ok := Find([]string{"storage-management", "update"})
	if !ok {
		t.Fatal("expected storage-management update command")
	}
	bodyJSON := `{"retention":"30d","version":3,"cold_move_duration":"7d","cold_volume":"cold","custom_rules":[]}`
	cfg := config.Config{
		APIKey:     "my-api-key",
		BackendID:  "my-backend",
		TenantUUID: "my-tenant",
		BaseURL:    srv.URL,
		Timeout:    time.Second,
	}
	var out bytes.Buffer
	err := Run(command, cfg, Options{
		PathValues: map[string]string{"dataType": "logs"},
		BodyJSON:   bodyJSON,
	}, &out)
	if err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
	if captured == nil {
		t.Fatal("server was never hit")
	}
	if readErr != nil {
		t.Fatalf("reading request body failed: %v", readErr)
	}
	if got, want := captured.Method, http.MethodPut; got != want {
		t.Errorf("method = %q, want %q", got, want)
	}
	if got, want := captured.URL.Path, "/api/storage-management/logs"; got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
	if got, want := captured.Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
	if got, want := captured.Header.Get("Authorization"), "Bearer my-api-key"; got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
	if got, want := captured.Header.Get("X-Tenant-UUID"), "my-tenant"; got != want {
		t.Errorf("X-Tenant-UUID = %q, want %q", got, want)
	}

	var gotBody, wantBody map[string]any
	if err := json.Unmarshal(capturedBody, &gotBody); err != nil {
		t.Fatalf("decoding captured body failed: %v", err)
	}
	if err := json.Unmarshal([]byte(bodyJSON), &wantBody); err != nil {
		t.Fatalf("decoding expected body failed: %v", err)
	}
	if !reflect.DeepEqual(gotBody, wantBody) {
		t.Errorf("body = %#v, want %#v", gotBody, wantBody)
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

// TestRunGrafanaWebApp_WithToken: grafana (WebApp) command with GrafanaToken set.
// Contract: Authorization is the Grafana bearer, X-Backend-Id and X-Tenant-UUID
// are absent even when TenantUUID is populated — grafana path bypasses SDK transport.
func TestRunGrafanaWebApp_WithToken(t *testing.T) {
	var captured *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	cmd := Command{
		Name:   []string{"grafana", "search"},
		Method: "GET",
		Path:   "/grafana/api/search",
		WebApp: true,
	}
	cfg := config.Config{
		GrafanaToken: "glsa_secret_token",
		TenantUUID:   "tenant-should-be-ignored",
		BaseURL:      srv.URL,
		Timeout:      time.Second,
	}
	var out bytes.Buffer
	if err := Run(cmd, cfg, Options{}, &out); err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
	if captured == nil {
		t.Fatal("server was never hit")
	}
	if got, want := captured.Header.Get("Authorization"), "Bearer glsa_secret_token"; got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
	if got := captured.Header.Get("X-Backend-Id"); got != "" {
		t.Errorf("X-Backend-Id = %q, want empty (grafana path must not inject backend-id)", got)
	}
	if got := captured.Header.Get("X-Tenant-UUID"); got != "" {
		t.Errorf("X-Tenant-UUID = %q, want empty (grafana path must not inject tenant-uuid)", got)
	}
	if out.Len() == 0 {
		t.Error("response body was not written to out")
	}
}

// TestRunGrafanaWebApp_MissingToken: grafana (WebApp) command with no GrafanaToken.
// Contract: Run returns an error mentioning the Grafana token and makes zero HTTP
// requests — the gcsa APIKey must not satisfy the grafana requirement.
func TestRunGrafanaWebApp_MissingToken(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cmd := Command{
		Name:   []string{"grafana", "search"},
		Method: "GET",
		Path:   "/grafana/api/search",
		WebApp: true,
	}
	cfg := config.Config{
		APIKey:  "gcsa-key-does-not-satisfy-grafana",
		BaseURL: srv.URL,
		Timeout: time.Second,
		// GrafanaToken intentionally empty
	}
	var out bytes.Buffer
	err := Run(cmd, cfg, Options{}, &out)
	if err == nil {
		t.Fatal("Run should return an error when GrafanaToken is missing")
	}
	lower := strings.ToLower(err.Error())
	if !strings.Contains(lower, "grafana") && !strings.Contains(lower, "service account") {
		t.Errorf("error %q should mention grafana or service account token", err.Error())
	}
	if hits != 0 {
		t.Errorf("server hit %d time(s), want 0 — error must short-circuit before HTTP request", hits)
	}
}

// TestRunNonWebApp_SDKHeaders: non-WebApp command with APIKey + BackendID + TenantUUID.
// Contract: SDK transport injects Authorization and X-Backend-Id; Run injects
// X-Tenant-UUID directly — all three reach the server.
func TestRunNonWebApp_SDKHeaders(t *testing.T) {
	var captured *http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	cmd := Command{
		Name:   []string{"backend", "settings"},
		Method: "GET",
		Path:   "/api/backend/settings",
		WebApp: false,
	}
	cfg := config.Config{
		APIKey:     "my-api-key",
		BackendID:  "my-backend",
		TenantUUID: "my-tenant",
		BaseURL:    srv.URL,
		Timeout:    time.Second,
	}
	var out bytes.Buffer
	if err := Run(cmd, cfg, Options{}, &out); err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
	if captured == nil {
		t.Fatal("server was never hit")
	}
	if got, want := captured.Header.Get("Authorization"), "Bearer my-api-key"; got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
	if got, want := captured.Header.Get("X-Backend-Id"), "my-backend"; got != want {
		t.Errorf("X-Backend-Id = %q, want %q", got, want)
	}
	if got, want := captured.Header.Get("X-Tenant-UUID"), "my-tenant"; got != want {
		t.Errorf("X-Tenant-UUID = %q, want %q", got, want)
	}
}
