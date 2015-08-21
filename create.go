package main

import (
	"encoding/base64"
	"fmt"
	"github.com/mattsurabian/msg"
	"github.com/spf13/cobra"
	"log"
)

var CreateCmd = &cobra.Command{
	Use:   "create [pact-name] [plain-text]",
	Short: "Outputs an encrypted ciphertext given a plain text message",
	Long: `Uses AES-256-GCM to encrypt a message with a randomly generated key from PBKDF2
and encrypts that secret key with the public key of each member of a pact. Base64 encoded encrypted
ciphertext is sent to STDOUT.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(Configuration.Pacts[args[0]]) > 0 {
			fmt.Println(Create(args[0], []byte(args[1])))
		} else {
			log.Fatalf("[ERROR] Config file does not contain keys for the pact: %s \n", args[0])
		}
	},
}

/**
 * Create
 * @param pactName string Name of the pact to encrypt data for
 * @param plainText []byte Byte array representation of the plain text
 * Uses the msg library to encrypt the provided plain text for consumption by the members of the
 * specified pact.
 */
func Create(pactName string, plainText []byte) string {
	pactKeyStrings := Configuration.Pacts[pactName]
	pactKeys := make([]*[32]byte, len(pactKeyStrings))

	for i, key := range pactKeyStrings {
		pactKeys[i] = msg.ReadNACLKeyString(key)
	}

	cipherText := msg.Encrypt(plainText, pactKeys, GetPublicKey(), GetPrivateKey())
	encodedCipherText := base64.StdEncoding.EncodeToString(cipherText)

	return encodedCipherText
}
