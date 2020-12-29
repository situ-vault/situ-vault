package vaultmode

import (
	"errors"
	"reflect"
)

type KeyDerivationFunction string

const (
	Pbkdf2_sha256_i10k    KeyDerivationFunction = "PBKDF2_SHA256_I10K"    // 10000 iterations, OpenSSL default (year 2000)
	Argon2id_t1_m65536_c4 KeyDerivationFunction = "ARGON2ID_T1_M65536_C4" // parameters as of RFC (year 2015)
	Scrypt_n32768_r8_p1   KeyDerivationFunction = "SCRYPT_N32768_R8_P1"   // parameters as RFC with bigger N (year 2009)
)

type keyDerivationFunctions struct {
	Pbkdf2_sha256_i10k    KeyDerivationFunction
	Argon2id_t1_m65536_c4 KeyDerivationFunction
	Scrypt_n32768_r8_p1   KeyDerivationFunction
}

var KeyDerivationFunctions = keyDerivationFunctions{
	Pbkdf2_sha256_i10k:    Pbkdf2_sha256_i10k,
	Argon2id_t1_m65536_c4: Argon2id_t1_m65536_c4,
	Scrypt_n32768_r8_p1:   Scrypt_n32768_r8_p1,
}

func ParseKeyDerivationFunction(s string) (KeyDerivationFunction, error) {
	for _, value := range allValues(reflect.ValueOf(KeyDerivationFunctions)) {
		if s == value {
			return KeyDerivationFunction(s), nil
		}
	}
	return "", errors.New("Invalid value: " + s)
}
