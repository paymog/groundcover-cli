package raw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/paymog/groundcover-cli/internal/config"
	"github.com/paymog/groundcover-cli/internal/output"
	"gopkg.in/yaml.v3"
)

const defaultWebAppBaseURL = "https://app.groundcover.com"

func Run(command Command, cfg config.Config, opts Options, out io.Writer) error {
	requestURL, err := buildURL(command, cfg, opts)
	if err != nil {
		return err
	}

	body, contentType, err := buildBody(command, opts)
	if err != nil {
		return err
	}

	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(command.Method, requestURL.String(), reader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json, text/event-stream, text/plain;q=0.9, */*;q=0.8")
	req.Header.Set("Accept-Encoding", "identity")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	var client *http.Client
	if command.WebApp {
		// Embedded Grafana is session-gated and ignores the gcsa bearer/backend
		// headers; it only accepts a Grafana service account token. Set that token
		// and bypass the SDK transport (which would clobber Authorization).
		if err := cfg.RequireGrafanaToken(); err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+cfg.GrafanaToken)
		client = cfg.WebAppHTTPClient()
	} else {
		if err := cfg.RequireAPIKey(); err != nil {
			return err
		}
		if cfg.TenantUUID != "" {
			req.Header.Set("X-Tenant-UUID", cfg.TenantUUID)
		}
		client = cfg.HTTPClient()
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d %s\n%s", resp.StatusCode, resp.Status, string(data))
	}
	return output.PrintBytes(out, data, cfg.Raw)
}

func buildURL(command Command, cfg config.Config, opts Options) (*url.URL, error) {
	path := command.Path
	for _, param := range command.PathParams {
		value := opts.PathValues[param]
		if value == "" {
			return nil, fmt.Errorf("missing --%s <%s>", kebab(param), param)
		}
		path = strings.ReplaceAll(path, ":"+param, url.PathEscape(value))
	}

	baseURL := cfg.NormalizedBaseURL()
	if command.WebApp && baseURL == config.DefaultBaseURL {
		baseURL = defaultWebAppBaseURL
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	rel := &url.URL{Path: path}
	requestURL := base.ResolveReference(rel)

	query := requestURL.Query()
	for key, value := range command.DefaultQuery {
		query.Set(key, value)
	}
	for _, item := range opts.Query {
		key, value, err := parseKeyValue(item)
		if err != nil {
			return nil, err
		}
		query.Set(key, value)
	}
	requestURL.RawQuery = query.Encode()
	return requestURL, nil
}

func buildBody(command Command, opts Options) ([]byte, string, error) {
	contentType := command.BodyContentType
	if contentType == "" {
		contentType = "application/json"
	}

	var body any
	if len(command.DefaultBody) > 0 {
		if err := json.Unmarshal(command.DefaultBody, &body); err != nil {
			return nil, "", err
		}
	}

	if opts.BodyFile != "" && opts.BodyJSON != "" {
		return nil, "", fmt.Errorf("use only one of --body-file or --body-json")
	}
	if opts.BodyFile != "" {
		data, err := os.ReadFile(opts.BodyFile)
		if err != nil {
			return nil, "", err
		}
		if isYAML(opts.BodyFile) {
			if len(opts.Set) == 0 {
				return data, defaultString(command.BodyContentType, "application/yaml"), nil
			}
			if err := yaml.Unmarshal(data, &body); err != nil {
				return nil, "", err
			}
		} else if err := json.Unmarshal(data, &body); err != nil {
			return nil, "", err
		}
	}
	if opts.BodyJSON != "" {
		if err := json.Unmarshal([]byte(opts.BodyJSON), &body); err != nil {
			return nil, "", err
		}
	}
	if len(opts.Set) > 0 {
		target, ok := body.(map[string]any)
		if !ok || target == nil {
			target = map[string]any{}
		}
		for _, item := range opts.Set {
			key, value, err := parseKeyValue(item)
			if err != nil {
				return nil, "", err
			}
			setDeep(target, key, parseValue(value))
		}
		body = target
	}

	if body == nil {
		return nil, "", nil
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, "", err
	}
	return data, contentType, nil
}

func parseKeyValue(input string) (string, string, error) {
	key, value, ok := strings.Cut(input, "=")
	if !ok || key == "" {
		return "", "", fmt.Errorf("expected key=value, got %q", input)
	}
	return key, value, nil
}

func parseValue(value string) any {
	switch value {
	case "true":
		return true
	case "false":
		return false
	case "null":
		return nil
	}
	var decoded any
	if err := json.Unmarshal([]byte(value), &decoded); err == nil {
		return decoded
	}
	return value
}

func setDeep(target map[string]any, path string, value any) {
	parts := strings.Split(path, ".")
	cursor := target
	for _, part := range parts[:len(parts)-1] {
		next, ok := cursor[part].(map[string]any)
		if !ok {
			next = map[string]any{}
			cursor[part] = next
		}
		cursor = next
	}
	cursor[parts[len(parts)-1]] = value
}

func isYAML(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}

func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
