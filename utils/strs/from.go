package strs

import (
	"encoding/json"

	"github.com/spf13/cast"
)

// From casts any value to a string type.
func From(v any) string {
	s, _ := FromE(v)
	return s
}

// From casts any value to a string type.
func FromE(v any) (string, error) {
	if d, ok := v.(json.RawMessage); ok {
		return string(d), nil
	}
	return cast.ToStringE(v)
}
