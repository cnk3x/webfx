package strs

import (
	"strings"

	"github.com/ettle/strcase"
)

func Snake[I BS](s I, upper ...bool) (out I) {
	var str string
	if len(upper) > 0 && upper[0] {
		str = strcase.ToSNAKE(string(s))
	} else {
		str = strcase.ToSnake(string(s))
	}
	str = strings.ReplaceAll(str, ".", "_")
	return I(str)
}

func Camel[I BS](s I) (out I) {
	return I(strcase.ToCamel(string(s)))
}

func Title[I BS](s I) (out I) {
	return I(strcase.ToCase(string(s), strcase.TitleCase, 0))
}
