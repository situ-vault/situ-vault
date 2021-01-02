package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/polarctos/situ-vault/pkg/vault"
)

func transform(manifest *Manifest, wd string) {
	envName := manifest.Config.Password.Env
	password, found := os.LookupEnv(envName)
	if !found {
		logStderr.Fatal("Failed to get password from environment: ", envName)
	}
	for _, filePath := range manifest.Files {
		transformFile(wd, filePath, password)
	}
}

func transformFile(wd string, filePath string, password string) {
	var content []byte
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		// also try with an explicit working directory:
		content, err = ioutil.ReadFile(wd + filePath)
		if err != nil {
			logStderr.Fatal("Failed to read file: ", os.Args[1])
		}
	}
	var secretManifest SecretManifest
	err = yaml.Unmarshal(content, &secretManifest)
	if err != nil {
		logStderr.Fatalf("Error unmarshalling manifest: %q \n%s\n", err, content)
	}
	transformed := content
	for _, value := range secretManifest.Data {
		cleartext, _, err := vault.Decrypt(value, password)
		if err != nil {
			logStderr.Fatalf("Decryption failed: %q \n%s\n", err, value)
		}
		encoded := base64.StdEncoding.EncodeToString([]byte(cleartext))
		transformed = bytes.Replace(transformed, []byte(value), []byte(encoded), 1)
	}
	// clear empty lines:
	clean := cleanManifest(string(transformed))
	logStdout.Println(clean)
	logStdout.Println("---")
}

func cleanManifest(manifest string) string {
	emptyLine := regexp.MustCompile(`(?m)\n\n`)
	s := emptyLine.ReplaceAllString(manifest, "\n")
	commentLine := regexp.MustCompile(`(?m)^# .*\n`)
	return commentLine.ReplaceAllString(s, "")
}

type SecretManifest struct {
	Kind string            `json:"kind" yaml:"kind"`
	Data map[string]string `json:"data" yaml:"data"`
}
