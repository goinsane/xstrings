package xstrings

import (
	"strconv"
	"unicode"
)

func TryUnquote(s string) string {
	str, err := strconv.Unquote(s)
	if err != nil {
		return s
	}
	return str
}

func ToUpperBeginning(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func ToLowerBeginning(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
