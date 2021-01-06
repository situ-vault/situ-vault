package internal

import (
	"runtime"
	"unicode"
)

var newlineChars = newlineOS()

func newlineOS() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func Wrap(text string, length int) string {
	var out string
	for i, r := range text {
		if i%length == 0 {
			out += newlineChars
		}
		out += string(r)
	}
	return out
}

func Unwrap(text string) string {
	var out string
	for _, r := range text {
		if !unicode.IsSpace(r) {
			out += string(r)
		}
	}
	return out
}
