package main

import (
	"encoding/json"
	"fmt"
	"go/format"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type har struct {
	Log struct {
		Entries []struct {
			Request struct {
				Method   string `json:"method"`
				URL      string `json:"url"`
				PostData *struct {
					Text     string `json:"text"`
					MimeType string `json:"mimeType"`
				} `json:"postData"`
			} `json:"request"`
			Response struct {
				Status int `json:"status"`
			} `json:"response"`
		} `json:"entries"`
	} `json:"log"`
}

type command struct {
	Name         []string
	Method       string
	Path         string
	Description  string
	PathParams   []string
	DefaultQuery map[string]string
	DefaultBody  []byte
	WebApp       bool
}

var (
	uuidPattern    = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	numericPattern = regexp.MustCompile(`^\d+$`)
	versionPattern = regexp.MustCompile(`^v\d+$`)
	excludedPaths  = map[string]bool{
		"/api/track/events":                        true,
		"/grafana/api/access-control/user/actions": true,
		"/grafana/api/frontend/assets":             true,
		"/grafana/api/frontend-metrics":            true,
	}
)

func main() {
	home, _ := os.UserHomeDir()
	harPath := filepath.Join(home, "Downloads", "app.groundcover.com.har")
	if len(os.Args) > 1 {
		harPath = os.Args[1]
	}

	data, err := os.ReadFile(harPath)
	must(err)

	var parsed har
	must(json.Unmarshal(data, &parsed))

	byEndpoint := map[string]command{}
	for _, entry := range parsed.Log.Entries {
		if entry.Request.URL == "" || entry.Request.Method == "" {
			continue
		}
		parsedURL, err := url.Parse(entry.Request.URL)
		if err != nil {
			continue
		}
		if !strings.HasSuffix(parsedURL.Hostname(), "groundcover.com") {
			continue
		}
		if !supportedPath(parsedURL.Path) {
			continue
		}
		if excludedPaths[parsedURL.Path] || entry.Response.Status >= 400 {
			continue
		}

		method := strings.ToUpper(entry.Request.Method)
		if !map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}[method] {
			continue
		}

		path, params := normalizedPath(parsedURL.Path)
		key := method + " " + path
		if _, exists := byEndpoint[key]; exists {
			continue
		}

		defaultQuery := map[string]string{}
		for key, values := range parsedURL.Query() {
			if len(values) > 0 {
				defaultQuery[key] = values[0]
			}
		}

		byEndpoint[key] = command{
			Name:         commandName(method, path),
			Method:       method,
			Path:         path,
			Description:  method + " " + path,
			PathParams:   params,
			DefaultQuery: defaultQuery,
			DefaultBody:  parseBody(entry.Request.PostData),
			WebApp:       strings.HasPrefix(parsedURL.Path, "/grafana/api/"),
		}
	}

	byName := map[string]command{}
	for _, command := range byEndpoint {
		key := strings.Join(command.Name, " ")
		byName[key] = prefer(command, byName[key])
	}

	commands := make([]command, 0, len(byName))
	for _, command := range byName {
		commands = append(commands, command)
	}
	sort.Slice(commands, func(i, j int) bool {
		return strings.Join(commands[i].Name, " ") < strings.Join(commands[j].Name, " ")
	})

	source, err := emit(commands)
	must(err)
	must(os.WriteFile("internal/raw/commands_generated.go", source, 0o644))
	fmt.Printf("Wrote %d raw commands to internal/raw/commands_generated.go\n", len(commands))
	if omitted := len(byEndpoint) - len(byName); omitted > 0 {
		fmt.Printf("Omitted %d older versioned endpoint aliases with duplicate versionless names\n", omitted)
	}
}

func supportedPath(path string) bool {
	return strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/grafana/api/")
}

func commandName(method string, path string) []string {
	name := []string{}
	trimmed := path
	if strings.HasPrefix(path, "/grafana/api/") {
		name = append(name, "grafana")
		trimmed = strings.TrimPrefix(path, "/grafana/api/")
	} else {
		trimmed = strings.TrimPrefix(path, "/api/")
	}

	parts := strings.Split(trimmed, "/")
	for i, part := range parts {
		if part == "" || versionPattern.MatchString(part) {
			continue
		}
		if (part == "uid" || part == "id") && i+1 < len(parts) && strings.HasPrefix(parts[i+1], ":") {
			continue
		}
		if strings.HasPrefix(part, ":") {
			if method == "GET" && i == len(parts)-1 && (len(name) == 0 || name[len(name)-1] != "get") {
				name = append(name, "get")
			}
			continue
		}
		name = append(name, part)
	}
	return name
}

func normalizedPath(path string) (string, []string) {
	parts := strings.Split(path, "/")
	params := []string{}
	for i, part := range parts {
		if !isDynamicPathPart(parts, i, part) {
			continue
		}
		name := dynamicParamName(parts, i)
		if !contains(params, name) {
			params = append(params, name)
		}
		parts[i] = ":" + name
	}
	return strings.Join(parts, "/"), params
}

func previousPart(parts []string, index int) string {
	for i := index - 1; i >= 0; i-- {
		if parts[i] != "" {
			return parts[i]
		}
	}
	return "id"
}

func isDynamicPathPart(parts []string, index int, part string) bool {
	if uuidPattern.MatchString(part) {
		return true
	}
	if !isGrafanaPath(parts) {
		return false
	}
	previous := previousPart(parts, index)
	switch previous {
	case "uid", "label":
		return true
	case "folders":
		return part != "folders"
	case "annotations", "id", "versions":
		return numericPattern.MatchString(part)
	default:
		return false
	}
}

func isGrafanaPath(parts []string) bool {
	return len(parts) > 3 && parts[1] == "grafana" && parts[2] == "api"
}

func dynamicParamName(parts []string, index int) string {
	previous := previousPart(parts, index)
	if isGrafanaPath(parts) {
		switch previous {
		case "uid":
			return uidParamName(previousPartBefore(parts, index-1))
		case "folders":
			return "folderUid"
		case "label":
			return "label"
		case "annotations":
			return "annotationId"
		case "id":
			return paramName(previousPartBefore(parts, index-1))
		case "versions":
			return "version"
		}
	}
	return paramName(previous)
}

func previousPartBefore(parts []string, index int) string {
	for i := index - 1; i >= 0; i-- {
		if parts[i] != "" {
			return parts[i]
		}
	}
	return "id"
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func paramName(previous string) string {
	base := strings.TrimSuffix(previous, "s")
	if base == "" {
		base = "id"
	}
	return base + "Id"
}

func uidParamName(previous string) string {
	base := strings.TrimSuffix(previous, "s")
	if base == "" {
		base = "id"
	}
	return base + "Uid"
}

func parseBody(postData *struct {
	Text     string `json:"text"`
	MimeType string `json:"mimeType"`
}) []byte {
	if postData == nil || postData.Text == "" {
		return nil
	}
	var decoded any
	if err := json.Unmarshal([]byte(postData.Text), &decoded); err != nil {
		encoded, _ := json.Marshal(postData.Text)
		return encoded
	}
	encoded, _ := json.Marshal(decoded)
	return encoded
}

func prefer(next command, current command) command {
	if current.Method == "" {
		return next
	}
	nextVersion := highestVersion(next.Path)
	currentVersion := highestVersion(current.Path)
	if nextVersion != currentVersion {
		if nextVersion > currentVersion {
			return next
		}
		return current
	}
	if len(next.Path) < len(current.Path) {
		return next
	}
	return current
}

func highestVersion(path string) int {
	highest := math.MinInt
	for _, part := range strings.Split(path, "/") {
		if !versionPattern.MatchString(part) {
			continue
		}
		version, _ := strconv.Atoi(strings.TrimPrefix(part, "v"))
		if version > highest {
			highest = version
		}
	}
	if highest == math.MinInt {
		return 0
	}
	return highest
}

func emit(commands []command) ([]byte, error) {
	var b strings.Builder
	b.WriteString("// Code generated by scripts/generate-commands.go from a Groundcover HAR; DO NOT EDIT.\n\n")
	b.WriteString("package raw\n\n")
	b.WriteString("import \"encoding/json\"\n\n")
	b.WriteString("var _ = json.RawMessage{}\n\n")
	b.WriteString("var Commands = []Command{\n")
	for _, command := range commands {
		b.WriteString("\t{\n")
		fmt.Fprintf(&b, "\t\tName: []string{%s},\n", quoteStrings(command.Name))
		fmt.Fprintf(&b, "\t\tMethod: %s,\n", strconv.Quote(command.Method))
		fmt.Fprintf(&b, "\t\tPath: %s,\n", strconv.Quote(command.Path))
		fmt.Fprintf(&b, "\t\tDescription: %s,\n", strconv.Quote(command.Description))
		if len(command.PathParams) > 0 {
			fmt.Fprintf(&b, "\t\tPathParams: []string{%s},\n", quoteStrings(command.PathParams))
		}
		if len(command.DefaultQuery) > 0 {
			fmt.Fprintf(&b, "\t\tDefaultQuery: map[string]string{%s},\n", quoteMap(command.DefaultQuery))
		}
		if len(command.DefaultBody) > 0 {
			fmt.Fprintf(&b, "\t\tDefaultBody: json.RawMessage(%s),\n", strconv.Quote(string(command.DefaultBody)))
		}
		if command.WebApp {
			b.WriteString("\t\tWebApp: true,\n")
		}
		b.WriteString("\t},\n")
	}
	b.WriteString("}\n")
	return format.Source([]byte(b.String()))
}

func quoteStrings(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, strconv.Quote(value))
	}
	return strings.Join(quoted, ", ")
}

func quoteMap(values map[string]string) string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, strconv.Quote(key)+": "+strconv.Quote(values[key]))
	}
	return strings.Join(parts, ", ")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
