package main

import (
	"github.com/mattsurabian/msg"
	"github.com/mitchellh/cli"

	"encoding/base64"
	"log"
	"strings"
)

type ReadPactCommand struct {
	UI cli.Ui
}

func (c *ReadPactCommand) Help() string {
	help := `
Usage: [flag] read [encrypted-message-base64]
 If no encrypted message is provided one will be read from standard input.
 This allows an encrypted and encoded message to be piped
 into the read command
`
	return strings.TrimSpace(help)
}

func (c *ReadPactCommand) Synopsis() string {
	return "Read an encrypted message"
}

func (c *ReadPactCommand) Run(args []string) int {

	encodedMessage := CheckStdIn()

	// Message has not been specified via std in, so it should be an argument
	if encodedMessage == "" && len(args) < 1 {
		log.Println("Error: Arguments missing. Run -h get for more info")
		return BAD_REQUEST
	}

	if encodedMessage == "" {
		encodedMessage = args[0]
	}

	message, err := base64.StdEncoding.DecodeString(encodedMessage)
	if err != nil {
		log.Println("Message Formatting Error (expected base64): ", err)
		return INTERNAL_ERROR
	}

	selfPublicKey := msg.ReadNACLKeyFile(GetNACLPublicKeyPath())
	selfPrivateKey := msg.ReadNACLKeyFile(GetNACLPrivateKeyPath())

	plainText, err := msg.Decrypt([]byte(message), selfPublicKey, selfPrivateKey)
	if err != nil {
		c.UI.Warn("Error decrypting message")
		return BAD_REQUEST
	}

	c.UI.Output(string(plainText))

	return OK
}
