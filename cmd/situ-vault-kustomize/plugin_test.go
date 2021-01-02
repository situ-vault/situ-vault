package main

import (
	"bytes"
	"encoding/base64"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	_ = os.Setenv("SITU_VAULT_PASSWORD", "test-pw")
	wd, _ := os.Getwd()

	os.Args = []string{"situ-vault-kustomize", "./testdata/example/secrets.yaml", wd + "/testdata/example/"}
	output, outputErr := captureOutput(main)

	assert.Equal(t, "", outputErr)
	assert.Contains(t, output, "kind: Secret\n")
	assert.Contains(t, output, "username: "+b64("test-data")+"\n")
	assert.Contains(t, output, "password: "+b64("test-data-longer")+"\n")

	assert.NotContains(t, "SITU_VAULT", output)
	assert.NotContains(t, ".yaml", output)
	assert.True(t, strings.HasSuffix(output, "---\n"))
}

// Helpers:

func b64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func captureOutput(function func()) (out string, err string) {
	var bufferOut bytes.Buffer
	var bufferErr bytes.Buffer
	logStdout.SetOutput(&bufferOut)
	logStderr.SetOutput(&bufferErr)
	function()
	return bufferOut.String(), bufferErr.String()
}
