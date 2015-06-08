package main

import (
	"github.com/mitchellh/cli"
	"os"
)

// Commands is the mapping of all available commands
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	Commands = map[string]cli.CommandFactory{
		"config": func() (cli.Command, error) {
			return &GenConfigCommand{
				UI: ui,
			}, nil
		},
		"key-gen": func() (cli.Command, error) {
			return &KeyGenCommand{
				UI: ui,
			}, nil
		},
		"read": func() (cli.Command, error) {
			return &ReadPactCommand{
				UI: ui,
			}, nil
		},
		"create": func() (cli.Command, error) {
			return &CreatePactCommand{
				UI: ui,
			}, nil
		},
		"key-export": func() (cli.Command, error) {
			return &KeyExportCommand{
				UI: ui,
			}, nil
		},
	}
}
