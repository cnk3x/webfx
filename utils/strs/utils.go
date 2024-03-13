package strs

import (
	"strings"
)

func PadLeft[I, P BS](s I, n int, pad P) I {
	return I(strings.Repeat(string(pad), n) + string(s))
}

func PadRight[I, P BS](s I, n int, pad P) I {
	return I(string(s) + strings.Repeat(string(pad), n))
}

func Clean[T comparable](c []T, fix ...func(item T) T) []T {
	var z T
	var x int
NEXT:
	for _, item := range c {
		v := item
		for _, f := range fix {
			if v = f(v); v == z {
				continue NEXT
			}
		}

		if v != z {
			c[x] = v
			x++
		}
	}

	return c[:x]
}
