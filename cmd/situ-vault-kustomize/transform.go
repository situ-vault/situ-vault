package main

import (
	"encoding/base64"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/situ-vault/situ-vault/pkg/vault"
)

func getPassword(passwordConfig PasswordConfig) string {
	var password string
	envName := passwordConfig.Env
	fileName := passwordConfig.File
	var found bool
	if envName != "" {
		password, found = os.LookupEnv(envName)
		if !found && fileName == "" {
			logStderr.Fatal("Failed to get password from environment: ", envName)
		}
	}
	if !found && fileName != "" {
		if strings.HasPrefix(fileName, "file:/") {
			password = fileName
			found = true
		} else {
			logStderr.Fatal("Password from file has to start with 'file:/'", envName)
		}
	}
	if !found {
		logStderr.Fatal("No password source defined")
	}
	return password
}

func transform(manifest *Manifest, wd string) {
	password := getPassword(manifest.Config.Password)
	for _, filePath := range manifest.Files {
		transformFile(wd, filePath, password)
	}
}

func transformFile(wd string, filePath string, password string) {
	var content []byte
	content, err := os.ReadFile(filePath)
	if err != nil {
		// also try with an explicit working directory:
		content, err = os.ReadFile(wd + filePath)
		if err != nil {
			logStderr.Fatal("Failed to read file: ", os.Args[1])
		}
	}
	var secretManifest SecretManifest
	err = yaml.Unmarshal(content, &secretManifest)
	if err != nil {
		logStderr.Fatalf("Error unmarshalling manifest: %q \n%s\n", err, content)
	}
	for key, value := range secretManifest.Data {
		cleartext, _, err := vault.Decrypt(value, password)
		if err != nil {
			logStderr.Fatalf("Decryption failed: %q \n%s\n", err, value)
		}
		encoded := base64.StdEncoding.EncodeToString([]byte(cleartext))
		secretManifest.Data[key] = encoded
	}
	transformed, err := yaml.Marshal(secretManifest)
	if err != nil {
		logStderr.Fatalf("Error marshalling transformed manifest: %q \n%s\n", err, content)
	}
	logStdout.Println(string(transformed))
	logStdout.Println("---")
}

type SecretManifest struct {
	ApiVersion string            `json:"apiVersion" yaml:"apiVersion,omitempty"`
	Kind       string            `json:"kind" yaml:"kind,omitempty"`
	Metadata   map[string]string `json:"metadata" yaml:"metadata,omitempty"`
	Type       string            `json:"type" yaml:"type,omitempty"`
	StringData map[string]string `json:"stringData" yaml:"stringData,omitempty"`
	Data       map[string]string `json:"data" yaml:"data,omitempty"`
}
