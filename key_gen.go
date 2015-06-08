package main

import (
	"github.com/mattsurabian/msg"
	"github.com/mitchellh/cli"
	"strings"
)

// KeyGenCommand generates and saves an NaCl public/private keypair
type KeyGenCommand struct {
	UI cli.Ui
}

// Long-form help
func (c *KeyGenCommand) Help() string {
	help := `
Usage: gen-key
  This command will save a NaCl public/private keypair
  to disk. NaCl is used for encryption and decryption
  of the random aes256-gcm keys
`
	return strings.TrimSpace(help)
}

func (c *KeyGenCommand) Synopsis() string {
	return "Create an public/private NaCl keypair"
}

// Run the actual command
func (c *KeyGenCommand) Run(args []string) int {
	if len(args) != 0 {
		c.UI.Info("The key-gen command does not accept arguments")
		c.UI.Info("Make sure to pass flags before the command")
		return BAD_REQUEST
	}

	if !DoesNACLKeypairExist() {
		pathErr := false
		privateKeyPath := GetNACLPrivateKeyPath()
		publicKeyPath := GetNACLPublicKeyPath()

		if privateKeyPath == "" {
			c.UI.Info("No configuration available for NaCl private key")
			c.UI.Info("Update your config file, or pass the -private-key flag to the client")
			pathErr = true
		}

		if publicKeyPath == "" {
			c.UI.Info("No configuration available for NaCl public key")
			c.UI.Info("Update your config file, or pass the -public-key flag to the client")
			pathErr = true
		}

		if pathErr {
			return BAD_REQUEST
		}

		c.UI.Info("Generating new keypair...")
		pubKey, privKey := msg.GenerateNACLKeyPair()
		msg.WriteNACLKeyFile(privateKeyPath, privKey, 0600)
		msg.WriteNACLKeyFile(publicKeyPath, pubKey, 0600)
	} else {
		c.UI.Info("A keypair already exists!")
		c.UI.Info("Change the config file to use different pub/priv key locations")
		c.UI.Info("Or pass -private-key and -public-key flags to the client")
	}
	return OK
}
