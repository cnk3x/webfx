package web

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/cnk3x/webfx/utils/strs"

	"github.com/Masterminds/sprig/v3"
)

var funcMap = sprig.TxtFuncMap()

func init() {
	funcMap["json"] = funcJSON
	funcMap["jmap"] = funcJmap
	funcMap["jarray"] = funcJarray
	funcMap["type"] = funcType
	funcMap["string"] = funcString
	funcMap["humanBytes"] = func(v reflect.Value) string {
		switch {
		case v.CanInt():
			return strs.HumanBytes(v.Int())
		case v.CanUint():
			return strs.HumanBytes(v.Uint())
		case v.CanFloat():
			return strs.HumanBytes(v.Float())
		case v.Kind() == reflect.String:
			return v.String()
		default:
			return ""
		}
	}
}

func funcJSON(v any, pretty ...bool) string {
	var data []byte
	if len(pretty) > 0 && pretty[0] {
		data, _ = json.MarshalIndent(v, "", "  ")
	} else {
		data, _ = json.Marshal(v)
	}
	return string(data)
}

func funcJmap(v json.RawMessage) (out map[string]any, err error) {
	if out = make(map[string]any); len(v) > 0 {
		err = json.Unmarshal([]byte(v), &out)
	}
	return
}

func funcJarray(v json.RawMessage) (out []map[string]any, err error) {
	if len(v) > 0 {
		err = json.Unmarshal([]byte(v), &out)
	}
	return
}

func funcType(v any) string { return fmt.Sprintf("%T", v) }

func funcString(v any) (s string) {
	var err error
	if s, err = strs.FromE(v); err != nil {
		s = fmt.Sprintf("%#v", v)
	}
	return
}
