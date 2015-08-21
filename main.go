package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var pactCmd = &cobra.Command{
	Use:   "pact",
	Short: "Pact is a CLI tool for encrypted communication between many parties.",
	Long: `A CLI tool that uses NaCl and AES-256-GCM to facilitate multiparty
communication without the need for out of band secret sharing.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Pact CLI v" + VERSION)
	},
}

func main() {
	pactCmd.AddCommand(CreateCmd)
	pactCmd.AddCommand(ReadCmd)
	pactCmd.AddCommand(ConfigCmd)
	pactCmd.AddCommand(KeyGenCmd)
	pactCmd.AddCommand(KeyExportCmd)
	pactCmd.AddCommand(NewPactCmd)
	pactCmd.AddCommand(DelPactCmd)
	pactCmd.AddCommand(ListPactCmd)
	pactCmd.AddCommand(AddPactKeyCmd)
	pactCmd.AddCommand(RmPactKeyCmd)
	pactCmd.Execute()
}
