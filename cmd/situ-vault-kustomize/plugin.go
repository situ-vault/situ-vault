package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Config Config   `json:"situVault" yaml:"situVault"`
	Files  []string `json:"files" yaml:"files"`
}
type Config struct {
	Password PasswordConfig `json:"password" yaml:"password"`
}
type PasswordConfig struct {
	Env  string `json:"env,omitempty" yaml:"env,omitempty"`
	File string `json:"file,omitempty" yaml:"file,omitempty"`
}

var logStdout = log.New(os.Stdout, "", 0)
var logStderr = log.New(os.Stderr, "", 0)

// kustomize exec plugin
// https://github.com/kubernetes-sigs/kustomize/blob/master/api/internal/plugins/execplugin/execplugin.go
func main() {
	if len(os.Args) < 2 {
		logStderr.Fatal("Wrong number of arguments:", os.Args)
	}

	// second argument is the file path
	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		logStderr.Fatal("Failed to read file: ", os.Args[1])
	}

	wd, err := os.Getwd()
	if err != nil || len(os.Args) > 2 {
		// overriding kustomize working directory, useful e.g. for tests
		wd = os.Args[2]
	}

	var manifest Manifest
	err = yaml.Unmarshal(content, &manifest)
	if err != nil {
		logStderr.Fatalf("Error unmarshalling manifest: %q \n%s\n", err, content)
	}

	transform(&manifest, wd)
}
