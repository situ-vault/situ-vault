package main

import (
	"flag"
	"log"
	"os"

	"github.com/situ-vault/situ-vault/pkg/vault"
	"github.com/situ-vault/situ-vault/pkg/vault/vaultmode"
)

var logStdout = log.New(os.Stdout, "", 0)
var logStderr = log.New(os.Stderr, "", 0)

func main() {
	result := handleCommand(os.Args)
	logStdout.Print(result) // careful this adds a newline at the end!
}

var defaultModeText = vaultmode.Defaults().Conservative.Text()

func handleCommand(args []string) string {
	encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	encryptPassword := encryptCmd.String("password", "", "the password")
	encryptCleartext := encryptCmd.String("cleartext", "", "the text to encrypt")
	encryptModeText := encryptCmd.String("vaultmode", defaultModeText, "the vaultmode")

	decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
	decryptPassword := decryptCmd.String("password", "", "the password")
	decryptCiphertext := decryptCmd.String("ciphertext", "", "the text to decrypt")

	if len(args) < 2 {
		logStderr.Fatal("expected 'encrypt' or 'decrypt' subcommands")
	}

	switch args[1] {
	case "encrypt":
		encryptCmd.Parse(args[2:])
		if *encryptPassword == "" {
			logStderr.Fatal("missing password")
		}
		if *encryptCleartext == "" {
			logStderr.Fatal("missing cleartext")
		}
		ciphertext, err := vault.Encrypt(*encryptCleartext, *encryptPassword, *encryptModeText)
		if err != nil {
			logStderr.Fatal("encrypt error: ", err)
		}
		return ciphertext
	case "decrypt":
		decryptCmd.Parse(args[2:])
		if *decryptPassword == "" {
			logStderr.Fatal("missing password")
		}
		if *decryptCiphertext == "" {
			logStderr.Fatal("missing ciphertext")
		}
		cleartext, _, err := vault.Decrypt(*decryptCiphertext, *decryptPassword)
		if err != nil {
			logStderr.Fatal("decrypt error: ", err)
		}
		return cleartext
	default:
		logStderr.Fatal("expected 'encrypt' or 'decrypt' subcommands")
		return "" // cannot be reached
	}
}
