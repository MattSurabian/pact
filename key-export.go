package main

import (
	"encoding/base64"
	"github.com/mattsurabian/msg"
	"github.com/mitchellh/cli"
	"strings"
)

// KeyExportCommand generates and saves an NaCl public/private keypair
type KeyExportCommand struct {
	UI cli.Ui
}

// Long-form help
func (c *KeyExportCommand) Help() string {
	help := `
Usage: export-key
  This command will write the contents of the currently loaded 
  NaCl public key to STDOUT with base64 encoding.
`
	return strings.TrimSpace(help)
}

func (c *KeyExportCommand) Synopsis() string {
	return "Write NaCl public key to STDOUT with base64 encoding"
}

// Run the actual command
func (c *KeyExportCommand) Run(args []string) int {
	if len(args) != 0 {
		c.UI.Info("The export-key command does not accept arguments")
		c.UI.Info("Make sure to pass flags before the command")
		return BAD_REQUEST
	}

	if DoesNACLKeypairExist() {
		publicKeyPath := GetNACLPublicKeyPath()

		if publicKeyPath == "" {
			c.UI.Info("No configuration value available for NaCl public key")
			c.UI.Info("Update your config file, or pass the -public-key flag to the client")
			return BAD_REQUEST
		}

		encodedKey := base64.StdEncoding.EncodeToString(msg.ReadNACLKeyFile(publicKeyPath)[:])
		c.UI.Info(encodedKey)
	} else {
		c.UI.Warn("Your public key cannot be found. Maybe generate one with the key-gen command?")
		return BAD_REQUEST
	}
	return OK
}
