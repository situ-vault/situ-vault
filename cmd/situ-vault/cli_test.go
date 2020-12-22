package main

import (
	"bytes"
	"github.com/polarctos/situ-vault/pkg/testdata"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var predefined = testdata.PredefinedDecrypt()

func Test_main(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"situ-vault", "encrypt", "-password=test-pw", "-cleartext=test-data"}
	outputEncrypt := captureOutput(main)
	assert.Contains(t, outputEncrypt, "SITU_VAULT")

	os.Args = []string{"situ-vault", "decrypt", "-password=" + predefined.Password, "-ciphertext=" + predefined.Ciphertext}
	outputDecrypt := captureOutput(main)
	assert.Contains(t, outputDecrypt, predefined.Cleartext)
}

func Test_handleCommand_roundTrip(t *testing.T) {
	password := testdata.RandomPassword(16)
	cleartext := testdata.RandomDataBase64(500)

	encryptArgs := []string{"situ-vault", "encrypt", "-password=" + password, "-cleartext=" + cleartext}
	resultEncrypted := handleCommand(encryptArgs)
	assert.Contains(t, resultEncrypted, "SITU_VAULT")
	assert.NotContains(t, resultEncrypted, cleartext)

	decryptArgs := []string{"situ-vault", "decrypt", "-password=" + password, "-ciphertext=" + resultEncrypted}
	resultDecrypted := handleCommand(decryptArgs)
	assert.Contains(t, resultDecrypted, cleartext)
}

func Test_handleCommand_decrypt_v1(t *testing.T) {
	decryptArgs := []string{"situ-vault", "decrypt", "-password=test-pw", "-ciphertext=SITU_VAULT_V1##AES256_GCM_PBKDF2_SHA256_ITER10K_SALT8_BASE32##TNSIVLVV6EOGI===##GRDENILPW24R4YDA2I6MKT6JPLG5GM2HWC5S2PR7"}
	resultDecrypted := handleCommand(decryptArgs)
	assert.Contains(t, resultDecrypted, "test-data")
}

// Helpers:

func captureOutput(function func()) string {
	var buffer bytes.Buffer
	log.SetOutput(&buffer)
	function()
	log.SetOutput(os.Stderr)
	return buffer.String()
}
