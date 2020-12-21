package vault

import (
	"errors"
	"strings"
)

const (
	prefix    = "SITU_VAULT_V1"
	mode      = "AES256_GCM_PBKDF2_SHA256_ITER10K_SALT8_BASE32"
	separator = "##"
)

func buildEnvelope(salt string, cipherText string) string {
	return prefix + separator + mode + separator + salt + separator + cipherText
}

func openEnvelope(data string) (salt string, cipherText string, err error) {
	split := strings.Split(data, separator)
	if split[0] != prefix {
		return "", "", errors.New("unknown input")
	}
	if split[1] != mode {
		return "", "", errors.New("unknown mode")
	}
	return split[2], split[3], nil
}
