package body

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Input struct {
	File string
	JSON string
}

func Decode(input Input, dest any) error {
	if input.File != "" && input.JSON != "" {
		return fmt.Errorf("use only one of --body-file or --body-json")
	}
	if input.JSON != "" {
		return json.Unmarshal([]byte(input.JSON), dest)
	}
	if input.File == "" {
		return fmt.Errorf("missing request body: pass --body-file or --body-json")
	}

	data, err := os.ReadFile(input.File)
	if err != nil {
		return err
	}
	if isYAML(input.File) {
		return yaml.Unmarshal(data, dest)
	}
	return json.Unmarshal(data, dest)
}

func isYAML(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}
