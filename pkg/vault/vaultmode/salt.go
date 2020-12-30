package vaultmode

import (
	"errors"
	"reflect"
)

type Salt string

const (
	R8b  Salt = "R8B"  // Random 8 bytes
	R16b Salt = "R16B" // Random 16 bytes
)

type salts struct {
	R8b  Salt
	R16b Salt
}

var Salts = salts{
	R8b:  R8b,
	R16b: R16b,
}

func ParseSalt(s string) (Salt, error) {
	for _, value := range allSaltValues {
		if s == value {
			return Salt(s), nil
		}
	}
	return "", errors.New("Invalid value: " + s)
}

var allSaltValues = allValues(reflect.ValueOf(Salts))

func (s salts) AllValues() []string {
	return allSaltValues
}
