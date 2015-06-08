package main

import (
	"github.com/mitchellh/cli"
	"log"
	"os"
)

func main() {
	c := cli.NewCLI("pact", VERSION)
	c.Args = os.Args[1:]
	c.Commands = Commands

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
