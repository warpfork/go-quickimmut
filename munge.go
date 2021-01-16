package main

import (
	"unicode"
)

func nostar(s string) string {
	switch {
	case s == "":
		return ""
	case s[0] == '*':
		return nostar(s)
	default:
		return s
	}
}

func upper(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(s)
	rs[0] = unicode.ToTitle(rs[0])
	return string(rs)
}

func lower(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(s)
	rs[0] = unicode.ToLower(rs[0])
	return string(rs)
}
