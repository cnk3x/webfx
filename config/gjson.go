package config

import (
	"github.com/tidwall/gjson"
)

func ParseJSON[T ~string | ~[]byte](data T) Value {
	return gjson.Parse(string(data))
}

type (
	Value = gjson.Result
	Type  = gjson.Type
)

const (
	Null   = gjson.Null
	False  = gjson.False
	Number = gjson.Number
	String = gjson.String
	True   = gjson.True
	JSON   = gjson.JSON
)
