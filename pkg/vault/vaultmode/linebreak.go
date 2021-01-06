package vaultmode

import (
	"errors"
	"reflect"
)

type Linebreak string

const (
	No    Linebreak = "NO"    // no line break
	Ch80  Linebreak = "CH80"  // after 80 characters
	Ch100 Linebreak = "CH100" // after 100 characters
	Ch120 Linebreak = "CH120" // after 120 characters
)

type linebreaks struct {
	No    Linebreak
	Ch80  Linebreak
	Ch100 Linebreak
	Ch120 Linebreak
}

var Linebreaks = linebreaks{
	No:    No,
	Ch80:  Ch80,
	Ch100: Ch100,
	Ch120: Ch120,
}

func ParseLinebreak(s string) (Linebreak, error) {
	for _, value := range allLinebreakValues {
		if s == value {
			return Linebreak(s), nil
		}
	}
	return "", errors.New("Invalid value: " + s)
}

var allLinebreakValues = allValues(reflect.ValueOf(Linebreaks))

func (s linebreaks) AllValues() []string {
	return allLinebreakValues
}
