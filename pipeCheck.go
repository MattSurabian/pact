package main

import (
	termutil "github.com/andrew-d/go-termutil"

	"io/ioutil"
	"log"
	"os"
)

// Figure out if we have something being piped into the CLI
// If so read and return it, if not return the empty string
func CheckStdIn() string {
	var val string
	if !termutil.Isatty(os.Stdin.Fd()) {
		valByteArr, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		val = string(valByteArr)
		return val
	}

	return val
}
