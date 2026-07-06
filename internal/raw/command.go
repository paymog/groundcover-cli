package raw

import (
	"encoding/json"
	"sort"
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
	WebApp          bool
}

func (c Command) Key() string {
	return strings.Join(c.Name, " ")
}

func Find(tokens []string) (Command, bool) {
	for _, command := range allCommands() {
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

func allCommands() []Command {
	commands := make([]Command, 0, len(Commands)+len(ExtraCommands))
	commands = append(commands, Commands...)
	commands = append(commands, ExtraCommands...)
	byName := map[string]Command{}
	for _, command := range commands {
		if _, exists := byName[command.Key()]; exists {
			continue
		}
		byName[command.Key()] = command
	}

	keys := make([]string, 0, len(byName))
	for key := range byName {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	ordered := make([]Command, 0, len(keys))
	for _, key := range keys {
		ordered = append(ordered, byName[key])
	}
	return ordered
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
