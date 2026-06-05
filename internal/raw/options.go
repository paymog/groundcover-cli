package raw

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/paymog/groundcover-cli/internal/config"
)

type Options struct {
	BodyFile   string
	BodyJSON   string
	Query      []string
	Set        []string
	PathValues map[string]string
}

func splitArgs(args []string) ([]string, []string) {
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			return args[:i], args[i:]
		}
	}
	return args, nil
}

func parseOptions(args []string, command Command, cfg *config.Config) (Options, error) {
	opts := Options{PathValues: map[string]string{}}
	pathFlags := map[string]string{}
	for _, param := range command.PathParams {
		pathFlags["--"+kebab(param)] = param
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		name, inlineValue, hasInlineValue := strings.Cut(arg, "=")

		if param, ok := pathFlags[name]; ok {
			value, err := optionValue(name, inlineValue, hasInlineValue, args, &i)
			if err != nil {
				return opts, err
			}
			opts.PathValues[param] = value
			continue
		}

		switch name {
		case "--body-file":
			value, err := optionValue(name, inlineValue, hasInlineValue, args, &i)
			if err != nil {
				return opts, err
			}
			opts.BodyFile = value
		case "--body-json":
			value, err := optionValue(name, inlineValue, hasInlineValue, args, &i)
			if err != nil {
				return opts, err
			}
			opts.BodyJSON = value
		case "--query":
			value, err := optionValue(name, inlineValue, hasInlineValue, args, &i)
			if err != nil {
				return opts, err
			}
			opts.Query = append(opts.Query, value)
		case "--set":
			value, err := optionValue(name, inlineValue, hasInlineValue, args, &i)
			if err != nil {
				return opts, err
			}
			opts.Set = append(opts.Set, value)
		case "--api-key":
			value, err := optionValue(name, inlineValue, hasInlineValue, args, &i)
			if err != nil {
				return opts, err
			}
			cfg.APIKey = value
		case "--backend-id":
			value, err := optionValue(name, inlineValue, hasInlineValue, args, &i)
			if err != nil {
				return opts, err
			}
			cfg.BackendID = value
		case "--tenant-uuid":
			value, err := optionValue(name, inlineValue, hasInlineValue, args, &i)
			if err != nil {
				return opts, err
			}
			cfg.TenantUUID = value
		case "--base-url":
			value, err := optionValue(name, inlineValue, hasInlineValue, args, &i)
			if err != nil {
				return opts, err
			}
			cfg.BaseURL = value
		case "--timeout":
			value, err := optionValue(name, inlineValue, hasInlineValue, args, &i)
			if err != nil {
				return opts, err
			}
			timeout, err := time.ParseDuration(value)
			if err != nil {
				return opts, fmt.Errorf("invalid --timeout: %w", err)
			}
			cfg.Timeout = timeout
		case "--raw":
			if hasInlineValue {
				raw, err := strconv.ParseBool(inlineValue)
				if err != nil {
					return opts, fmt.Errorf("invalid --raw value: %w", err)
				}
				cfg.Raw = raw
			} else {
				cfg.Raw = true
			}
		default:
			return opts, fmt.Errorf("unknown raw option %s", name)
		}
	}

	return opts, nil
}

func optionValue(name string, inlineValue string, hasInlineValue bool, args []string, index *int) (string, error) {
	if hasInlineValue {
		return inlineValue, nil
	}
	*index++
	if *index >= len(args) || strings.HasPrefix(args[*index], "-") {
		return "", fmt.Errorf("missing value for %s", name)
	}
	return args[*index], nil
}
