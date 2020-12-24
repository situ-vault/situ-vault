package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	_ = os.Setenv("SITU_VAULT_PASSWORD", "test-pw")
	wd, _ := os.Getwd()

	os.Args = []string{"situ-vault-kustomize", "./testdata/example/secrets.yaml", wd + "/testdata/example/"}
	output := captureOutput(main)

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

func captureOutput(function func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()
	os.Stdout = writer
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		written, err2 := io.Copy(&buf, reader)
		if err2 != nil || written == 0 {
			panic("Failed to copy output")
		}
		out <- buf.String()
	}()
	wg.Wait()
	function()
	err = writer.Close()
	if err != nil {
		panic(err)
	}
	return <-out
}
