package output

import (
	"encoding/json"
	"fmt"
	"io"
)

func Print(w io.Writer, value any, raw bool) error {
	switch v := value.(type) {
	case nil:
		_, err := fmt.Fprintln(w, "null")
		return err
	case []byte:
		return PrintBytes(w, v, raw)
	case string:
		return PrintBytes(w, []byte(v), raw)
	default:
		data, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, string(data))
		return err
	}
}

func PrintBytes(w io.Writer, data []byte, raw bool) error {
	if raw {
		_, err := w.Write(data)
		return err
	}

	var decoded any
	if err := json.Unmarshal(data, &decoded); err == nil {
		return Print(w, decoded, false)
	}
	_, err := w.Write(data)
	return err
}
