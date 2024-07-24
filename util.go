package main

import (
	"strings"
	"unicode"
)

func StripStr(s string) string {
	return strings.Map(func(r rune) rune {
		lower := unicode.ToLower(r)
		if lower >= 'a' && lower <= 'z' {
			return lower
		}
		return -1
	}, s)
}
