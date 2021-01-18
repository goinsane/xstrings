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
	buf := make([]byte, 0, len(str))
	var idx int
	for _, r := range str {
		if idx <= 0 {
			r = unicode.ToUpper(r)
		}
		buf = append(buf, []byte(string(r))...)
		idx++
	}
	return string(buf)
}

func ToLowerBeginning(str string) string {
	if str == "" {
		return ""
	}
	buf := make([]byte, 0, len(str))
	var idx int
	for _, r := range str {
		if idx <= 0 {
			r = unicode.ToLower(r)
		}
		buf = append(buf, []byte(string(r))...)
		idx++
	}
	return string(buf)
}
