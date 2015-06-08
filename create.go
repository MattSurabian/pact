package main

import (
	"github.com/mattsurabian/msg"
	"github.com/mitchellh/cli"

	"encoding/base64"
	"log"
	"strings"
)

// CreatePactCommand creates an encrypted message
type CreatePactCommand struct {
	UI cli.Ui
}

// Long-form help
func (c *CreatePactCommand) Help() string {
	help := `
Usage: [flags] create [message]
Output: A base64 encoded representation of the encrypted message.
 If no message is provided standard input will be read.
 This allows text to be piped into the create command
`
	return strings.TrimSpace(help)
}

func (c *CreatePactCommand) Synopsis() string {
	return "Create an encrypted message"
}

// Run the actual command
func (c *CreatePactCommand) Run(args []string) int {

	ptMessage := []byte(CheckStdIn())

	// Message has not been specified via std in, so it should be an argument
	if len(ptMessage) == 0 && len(args) < 1 {
		log.Println("Error: Missing arguments, run -h put for more info")
		return BAD_REQUEST
	}

	if len(ptMessage) == 0 {
		ptMessage = []byte(args[0])
	}

	selfPublicKey := msg.ReadNACLKeyFile(GetNACLPublicKeyPath())
	selfPrivateKey := msg.ReadNACLKeyFile(GetNACLPrivateKeyPath())

	// In future this should be a call to the server which returns a byte array of pub keys
	authorizedUserKeys := []*[32]byte{selfPublicKey}

	cipherText := msg.Encrypt(ptMessage, authorizedUserKeys, selfPublicKey, selfPrivateKey)
	encodedCipherText := base64.StdEncoding.EncodeToString(cipherText)

	c.UI.Output(encodedCipherText)

	return OK
}
