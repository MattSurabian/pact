package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/mattsurabian/msg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var KeyGenCmd = &cobra.Command{
	Use:   "key-gen",
	Short: "Creates new NaCL keys in the location specified by pact's configuration",
	Long:  `Generates an NaCL keypair and writes their base64 string representation to the files specified in pact's current configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		KeyGen()
	},
}

var KeyExportCmd = &cobra.Command{
	Use:   "key-export",
	Short: "Outputs the user's public key encoded as base64 to STDOUT",
	Long:  `Sends the user's public key encoded as base64 to STDOUT for easy distribution`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(base64.StdEncoding.EncodeToString(GetPublicKey()[:]))
	},
}

func init() {
	viper.SetDefault("PublicKeyPath", ConfigDirectory+"naclPub.key")
	viper.SetDefault("PrivateKeyPath", ConfigDirectory+"naclPriv.key")
}

/*
 * KeyGen
 * @returns err error Returns nil if no errors are encountered
 * Method which generates a new keypair at the locations defined in the config. If one or both
 * pieces of the keypair already exist this method returns an error.
 */
func KeyGen() (err error) {
	pubKeyPath := GetPublicKeyPath()
	privKeyPath := GetPrivateKeyPath()

	if !DoesNACLKeypairExist() {
		fmt.Println("Generating new keypair...")
		pubKey, privKey := msg.GenerateNACLKeyPair()
		msg.WriteNACLKeyFile(pubKeyPath, pubKey, 0600)
		msg.WriteNACLKeyFile(privKeyPath, privKey, 0600)
		fmt.Println("Keypair created! ")

	} else {
		err = errors.New("Keypair already exists!")
		fmt.Println("Parts of your keypair already exist! Refusing to overwrite. Is your configuration correct?")
	}
	fmt.Println(" PublicKeyPath: " + pubKeyPath)
	fmt.Println(" PrivateKeyPath: " + privKeyPath)
	return
}

/**
 * DoesNACLKeypairExist
 * @returns bool
 * Helper method which checks for the existence of public and private keys based on
 * the currently loaded configurationManager values.
 */
func DoesNACLKeypairExist() bool {
	_, pubKeyError := os.Stat(GetPublicKeyPath())
	_, privKeyError := os.Stat(GetPrivateKeyPath())
	return !(os.IsNotExist(pubKeyError) && os.IsNotExist(privKeyError))
}

/**
 * GetPublicKeyAbsPath
 * @returns string Path to the public key
 * Helper method that returns the path to the public key
 */
func GetPublicKeyPath() string {
	return viper.GetString("PublicKeyPath")
}

/**
 * GetPublicKey
 * @return *[32]byte Byte array representation of the user's public key
 * Helper method which returns the user's public key
 */
func GetPublicKey() *[32]byte {
	return msg.ReadNACLKeyFile(GetPublicKeyPath())
}

/**
 * GetPrivateKeyPath
 * @returns string Path to the private key
 * Helper method that returns the path to the private key
 */
func GetPrivateKeyPath() string {
	return viper.GetString("PrivateKeyPath")
}

/**
 * GetPrivateKey
 * @return *[32]byte Byte array representation of the user's private key
 * Helper method which returns the user's private key
 */
func GetPrivateKey() *[32]byte {
	return msg.ReadNACLKeyFile(GetPrivateKeyPath())
}
