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

func IsBeginningUpper(str string) bool {
	if str == "" {
		return false
	}
	return unicode.IsUpper([]rune(str)[0])
}

func IsBeginningLower(str string) bool {
	if str == "" {
		return false
	}
	return unicode.IsLower([]rune(str)[0])
}

func AreLettersUpper(str string) bool {
	for _, r := range []rune(str) {
		if !unicode.IsLetter(r) {
			continue
		}
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func AreLettersLower(str string) bool {
	for _, r := range []rune(str) {
		if !unicode.IsLetter(r) {
			continue
		}
		if !unicode.IsLower(r) {
			return false
		}
	}
	return true
}
