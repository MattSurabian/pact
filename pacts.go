package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var NewPactCmd = &cobra.Command{
	Use:   "new [pact-name]",
	Short: "Creates a new pact",
	Long:  `Creates a new pact in the configuration file that keys can be added to with the add-key command`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing pact name! Refusing to continue")
			os.Exit(400)
		}
		pactName := args[0]
		Configuration.Pacts[pactName] = []string{}
		PersistConfiguration()
	},
}

var ListPactCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists existing pacts",
	Long:  `Outputs a list of existing pacts and the keys they contain.`,
	Run: func(cmd *cobra.Command, args []string) {
		pacts, err := json.MarshalIndent(Configuration.Pacts, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(pacts))
	},
}

var AddPactKeyCmd = &cobra.Command{
	Use:   "add-key [pact-name] [public-key]",
	Short: "Adds a key to an existing pact or creates a new pact containing the key",
	Long: `Adds the provided public key to the specified pact. A new pact will be
created if necessary. The public-key can be piped into this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing pact name! Refusing to continue")
			os.Exit(400)
		}
		pactName := args[0]

		key := CheckStdIn()
		if key == "" && len(args) < 2 {
			fmt.Println("Missing public key to add! It can be specified as the second parameter or piped into the command. Refusing to continue")
			os.Exit(400)
		}

		if key == "" {
			key = args[1]
		}

		Configuration.Pacts[pactName] = append(Configuration.Pacts[pactName], strings.TrimSpace(key))
		PersistConfiguration()
	},
}

var RmPactKeyCmd = &cobra.Command{
	Use:   "rm-key [pact-name]",
	Short: "Interactively removes a single key from an existing pact",
	Long:  `Removes a single key from an existing pact using interactive prompts.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing pact name! Refusing to continue")
			os.Exit(400)
		}
		pactName := args[0]
		pactKeyStrings := Configuration.Pacts[pactName]

		if len(pactKeyStrings) == 0 {
			fmt.Println("No such pact: " + pactName)
			os.Exit(400)
		}

		for i, key := range pactKeyStrings {
			fmt.Println("[" + strconv.Itoa(i) + "] " + key)
		}

		var removalIndex int
		fmt.Println("Which key would you like to remove?")
		_, err := fmt.Scanf("%d", &removalIndex)
		if err != nil || removalIndex > len(pactKeyStrings) {
			fmt.Println("Unexpected input, exiting...")
			os.Exit(400)
		}

		Configuration.Pacts[pactName] = append(pactKeyStrings[:removalIndex], pactKeyStrings[removalIndex+1:]...)

		PersistConfiguration()
	},
}

var DelPactCmd = &cobra.Command{
	Use:   "rm [pact-name]",
	Short: "Completely removes an existing pact",
	Long:  `Removes an existing pact and all the keys it contains from the user's configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing pact name! Refusing to continue")
			os.Exit(400)
		}
		pactName := args[0]
		delete(Configuration.Pacts, pactName)
		PersistConfiguration()
	},
}

/**
 * SetDefaultPact
 * Method which sets the "self" pact containing the user's public key. If we can't read
 * THIS configuration model's public key THIS configuration model should have an empty "self"
 * pact. In that way we avoid assuming that the user has already generated a key in the provided
 * location and can attempt to help the user by automatically generating a new keypair by calling
 * KeyGen and then recursively calling this function to ensure the persisted configuration file
 * always contains a "self" pact.
 *
 * The safety net here is that if part of the keypair is found, a private key without the public or
 * vice-versa, KeyGen will error, no keys will be generated, and no recursive call of
 * this method will occur.
 */
func SetDefaultPact() {
	selfPublicKey, err := ioutil.ReadFile(GetPublicKeyPath())
	if err == nil {
		viper.Set("Pacts", map[string][]string{"self": []string{string(selfPublicKey)}})
	} else {
		viper.Set("Pacts", map[string][]string{})
		fmt.Println("Your keypair doesn't exist, automatically generating one...")
		keyGenError := KeyGen()
		if keyGenError == nil {
			SetDefaultPact()
		}
	}
}
