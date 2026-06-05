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
}

var (
	uuidPattern    = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	versionPattern = regexp.MustCompile(`^v\d+$`)
	excludedPaths  = map[string]bool{"/api/track/events": true}
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
		if !strings.HasPrefix(parsedURL.Path, "/api/") {
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
			Name:         commandName(parsedURL.Path),
			Method:       method,
			Path:         path,
			Description:  method + " " + path,
			PathParams:   params,
			DefaultQuery: defaultQuery,
			DefaultBody:  parseBody(entry.Request.PostData),
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

func commandName(path string) []string {
	parts := strings.Split(strings.TrimPrefix(path, "/api/"), "/")
	name := []string{}
	for i, part := range parts {
		if part == "" || versionPattern.MatchString(part) {
			continue
		}
		if uuidPattern.MatchString(part) {
			if i == 0 || len(name) == 0 || name[len(name)-1] != "get" {
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
		if !uuidPattern.MatchString(part) {
			continue
		}
		name := paramName(previousPart(parts, i))
		params = append(params, name)
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

func paramName(previous string) string {
	base := strings.TrimSuffix(previous, "s")
	if base == "" {
		base = "id"
	}
	return base + "Id"
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
