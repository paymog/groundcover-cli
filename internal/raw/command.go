package raw

import (
	"encoding/json"
	"strings"
)

type Command struct {
	Name            []string
	Method          string
	Path            string
	Description     string
	PathParams      []string
	DefaultQuery    map[string]string
	DefaultBody     json.RawMessage
	BodyContentType string
}

func (c Command) Key() string {
	return strings.Join(c.Name, " ")
}

func Find(tokens []string) (Command, bool) {
	for _, command := range Commands {
		if len(command.Name) != len(tokens) {
			continue
		}
		matches := true
		for i := range command.Name {
			if command.Name[i] != tokens[i] {
				matches = false
				break
			}
		}
		if matches {
			return command, true
		}
	}
	return Command{}, false
}

func kebab(value string) string {
	var out strings.Builder
	for i, r := range value {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				out.WriteByte('-')
			}
			out.WriteRune(r + ('a' - 'A'))
			continue
		}
		out.WriteRune(r)
	}
	return out.String()
}
