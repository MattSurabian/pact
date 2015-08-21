package main

import (
	"encoding/base64"
	"fmt"
	"github.com/mattsurabian/msg"
	"github.com/spf13/cobra"
	"log"
)

var ReadCmd = &cobra.Command{
	Use:   "read [ciphertext]",
	Short: "Outputs a plain text message given an encrypted ciphertext",
	Long: `Uses NaCl to decrypt a key which can be used to decrypt the message
which has been secured with AES-256-GCM encryption.`,
	Run: func(cmd *cobra.Command, args []string) {
		encodedMessage := CheckStdIn()

		// Message has not been specified via std in, so it should be an argument
		if encodedMessage == "" && len(args) < 1 {
			log.Fatalln("[ERROR] Arguments missing. Run -h get for more info")
		}

		if encodedMessage == "" {
			encodedMessage = args[0]
		}

		message, err := base64.StdEncoding.DecodeString(encodedMessage)
		if err != nil {
			log.Fatalln("[ERROR] Message Formatting Error (expected base64)")
		}

		fmt.Println(Read(message))
	},
}

/**
 * Read
 * @param cipherText []byte Ciphertext to decrypt
 * @returns string Plain text of the decrypted ciphertext
 * Uses the msg library to decrypt a byte array of ciphertext and returns the plain text.
 * On error bails on execution.
 */
func Read(cipherText []byte) string {
	plainText, err := msg.Decrypt(cipherText, GetPublicKey(), GetPrivateKey())
	if err != nil {
		log.Fatalln("[ERROR] Decryption Failed")
	}

	return string(plainText)
}
