package mode

import (
	"fmt"
	"reflect"
	"strings"
)

// examples: (the mode is only the part following after "SITU_VAULT_V1##" until the next "##")
// "SITU_VAULT_V1##C:AES256_GCM#KDF:PBKDF2_SHA256_I10K#SALT:R8B#ENC:BASE32##TNSIVLVV6EOGI===##GRDENILPW24R4YDA2I6MKT6JPLG5GM2HWC5S2PR7"
// WIP: "SITU_VAULT_V1##C:AES256_GCM#KDF:ARGON2ID_T1_M65536_C4#SALT:R8B#ENC:BASE64URL##YW55IGNhc===##jA0EAw_BBPOQhPfTDInn-94hXmnBr9D8-4x5"
// WIP: "SITU_VAULT_V1##C:NACL_SECRETBOX#KDF:SCRYPT_N32768_R8_P1#SALT:R16B#ENC:BASE64##YW55IGNhcm5hbCBwbGV===##jA0EAw/BBPOQhPfTDInn+94hXmnBr9D8+4x5"

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

type Construct string

const (
	AES256_GCM         Construct = "AES256_GCM"         // AEAD; Standard: Nonce: 12 byte, Tag: 16 byte (year 2000/2007)
	NACL_SECRETBOX     Construct = "NACL_SECRETBOX"     // AE; XSalsa20-Poly1305; Standard: Nonce: 24 byte, Tag: 16 byte (year 2008)
	XCHACHA20_POLY1305 Construct = "XCHACHA20_POLY1305" // AEAD; XChaCha20-Poly1305; Standard: Nonce: 24 byte, Tag: 16 byte (year 2008/2018)
)

type KeyDerivationFunction string

const (
	PBKDF2_SHA256_I10K    KeyDerivationFunction = "PBKDF2_SHA256_I10K"    // 10000 iterations, OpenSSL default (year 2000)
	ARGON2ID_T1_M65536_C4 KeyDerivationFunction = "ARGON2ID_T1_M65536_C4" // parameters as of RFC (year 2015)
	SCRYPT_N32768_R8_P1   KeyDerivationFunction = "SCRYPT_N32768_R8_P1"   // parameters as RFC with bigger N (year 2009)
)

type Salt string

const (
	R8B  Salt = "R8B"  // Random 8 bytes
	R16B Salt = "R16B" // Random 16 bytes
)

type Encoding string

const (
	NONE      Encoding = "NONE"      // No encoding, just bytes
	BASE32    Encoding = "BASE32"    // Base32
	BASE62    Encoding = "BASE62"    // Base62 (is Base64 without the 2 special characters)
	BASE64    Encoding = "BASE64"    // Base64
	BASE64URL Encoding = "BASE64URL" // Base64 (URL safe variant)
)

func (m Mode) Text() string {
	v := reflect.ValueOf(m)
	var text string
	for i := 0; i < v.NumField(); i++ {
		code, found := v.Type().Field(i).Tag.Lookup(code)
		if !found {
			panic("Untagged mode field: " + v.Type().Field(i).Name)
		}
		fieldValue := fmt.Sprint(v.Field(i).Interface())
		text += code + codeSeparator + fieldValue
		if i != v.NumField()-1 {
			text += fieldSeparator
		}
	}
	return text
}

func NewMode(text string) *Mode {
	split := strings.Split(text, fieldSeparator)
	m := &Mode{
		Construct: "",
		Kdf:       "",
		Salt:      "",
		Encoding:  "",
	}
	v := reflect.Indirect(reflect.ValueOf(m))
	if len(split) != v.NumField() {
		panic("Incorrect mode text: " + text)
	}
	for i := 0; i < v.NumField(); i++ {
		code, found := v.Type().Field(i).Tag.Lookup(code)
		structFieldName := v.Type().Field(i).Name
		if !found {
			panic("Untagged mode field: " + structFieldName)
		}
		parts := strings.Split(split[i], codeSeparator)
		fieldCode := parts[0]
		if code != fieldCode {
			panic("Incorrect mode order: " + text + " " + fieldCode)
		}
		fieldValue := parts[1]
		structField := v.Field(i)
		if !structField.CanSet() && structField.Kind() != reflect.String {
			panic("Cannot set field: " + structFieldName)
		}
		v.Field(i).SetString(fieldValue)
	}
	return m
}

func (m Mode) DeepEqual(m2 *Mode) bool {
	// reflect.DeepEqual() uses == to compare strings, this fits
	return reflect.DeepEqual(m, m2)
}
