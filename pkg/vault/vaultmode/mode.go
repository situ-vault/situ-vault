package vaultmode

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// examples: (the vaultmode is only the part following after "SITU_VAULT_V1##" until the next "##")
// "SITU_VAULT_V1##C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32##TNSIVLVV6EOGI===##GRDENILPW24R4YDA2I6MKT6JPLG5GM2HWC5S2PR7##END"

const (
	code           string = `code` // field tag used to get the type prefix for the text representation from the struct
	fieldSeparator string = `#`    // separator used in the text representation
	codeSeparator  string = ":"
)

type Mode struct {
	Construct Construct             `code:"C"`
	Kdf       KeyDerivationFunction `code:"KDF"`
	Salt      Salt                  `code:"SALT"`
	Encoding  Encoding              `code:"ENC"`
}

func (m Mode) Text() string {
	v := reflect.ValueOf(m)
	var text string
	for i := 0; i < v.NumField(); i++ {
		code, found := v.Type().Field(i).Tag.Lookup(code)
		if !found {
			panic("Untagged vaultmode field: " + v.Type().Field(i).Name)
		}
		fieldValue := fmt.Sprint(v.Field(i).Interface())
		text += code + codeSeparator + fieldValue
		if i != v.NumField()-1 {
			text += fieldSeparator
		}
	}
	return text
}

func NewMode(text string) (*Mode, error) {
	split := strings.Split(text, fieldSeparator)
	m := &Mode{
		Construct: "",
		Kdf:       "",
		Salt:      "",
		Encoding:  "",
	}
	v := reflect.Indirect(reflect.ValueOf(m))
	if len(split) != v.NumField() {
		return nil, errors.New("Incorrect vaultmode text: " + text)
	}
	for i := 0; i < v.NumField(); i++ {
		code, found := v.Type().Field(i).Tag.Lookup(code)
		structFieldName := v.Type().Field(i).Name
		if !found {
			return nil, errors.New("Untagged vaultmode field: " + structFieldName)
		}
		parts := strings.Split(split[i], codeSeparator)
		fieldCode := parts[0]
		if code != fieldCode {
			return nil, errors.New("Incorrect vaultmode order: " + text + " " + fieldCode)
		}
		fieldValue := parts[1]
		structField := v.Field(i)
		if !structField.CanSet() && structField.Kind() != reflect.String {
			return nil, errors.New("Cannot set field: " + structFieldName)
		}
		err := validateModeField(fieldCode, fieldValue)
		if err != nil {
			return nil, err
		}
		v.Field(i).SetString(fieldValue)
	}
	return m, nil
}

func validateModeField(code string, value string) (err error) {
	switch code {
	case "C":
		_, err = ParseConstruct(value)
	case "KDF":
		_, err = ParseKeyDerivationFunction(value)
	case "SALT":
		_, err = ParseSalt(value)
	case "ENC":
		_, err = ParseEncoding(value)
	}
	return err
}

func allValues(v reflect.Value) (e []string) {
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		e = append(e, value.String())
	}
	return e
}
