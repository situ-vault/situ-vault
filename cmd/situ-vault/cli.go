package main

import (
	"flag"
	"log"
	"os"

	"github.com/polarctos/situ-vault/pkg/vault"
	"github.com/polarctos/situ-vault/pkg/vault/mode"
)

func main() {
	log.Println("Running situ-vault")

	result := handleCommand(os.Args)
	log.Println("Result: " + result)
}

var defaultModeText string = mode.Defaults().Conservative.Text()

func handleCommand(args []string) string {
	encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	encryptPassword := encryptCmd.String("password", "", "the password")
	encryptCleartext := encryptCmd.String("cleartext", "", "the text to encrypt")
	encryptModeText := encryptCmd.String("mode", defaultModeText, "the mode")

	decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
	decryptPassword := decryptCmd.String("password", "", "the password")
	decryptCiphertext := decryptCmd.String("ciphertext", "", "the text to decrypt")

	if len(args) < 2 {
		log.Println("expected 'encrypt' or 'decrypt' subcommands")
		os.Exit(1)
	}

	switch args[1] {
	case "encrypt":
		encryptCmd.Parse(args[2:])
		ciphertext, err := vault.Encrypt(*encryptCleartext, *encryptPassword, *encryptModeText)
		if err != nil {
			log.Fatal("encrypt error: ", err)
		}
		return ciphertext
	case "decrypt":
		decryptCmd.Parse(args[2:])
		cleartext, _, err := vault.Decrypt(*decryptCiphertext, *decryptPassword)
		if err != nil {
			log.Fatal("decrypt error: ", err)
		}
		return cleartext
	default:
		log.Println("expected 'encrypt' or 'decrypt' subcommands")
		os.Exit(1)
		return "" // cannot be reached
	}
}
