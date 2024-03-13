package strs

import (
	"github.com/spf13/cast"
)

func ToInt[o Int](s any) o {
	return o(cast.ToInt64(s))
}

func ToIntOr[o Int](s any, def o) o {
	i, err := cast.ToInt64E(s)
	if err != nil {
		return def
	}
	return o(i)
}
