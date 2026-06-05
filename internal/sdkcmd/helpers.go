package sdkcmd

import "github.com/paymog/groundcover-cli/internal/body"

// decodeOptionalBody decodes a request body if one was supplied, otherwise leaves dest at its zero value.
// Used for endpoints whose body is a filter (e.g. {query: "..."}) and is optional.
func decodeOptionalBody(input body.Input, dest any) error {
	if input.File == "" && input.JSON == "" {
		return nil
	}
	return body.Decode(input, dest)
}
