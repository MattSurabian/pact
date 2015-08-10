package main

import (
	termutil "github.com/andrew-d/go-termutil"

	"io/ioutil"
	"log"
	"os"
)

/**
 * CheckStdIn
 * @returns string value from STDIN
 * Helper method which figures out if something is being piped into the CLI. If so, it is
 * read and returned. Otherwise this method returns the empty string.
 */
func CheckStdIn() string {
	var val string
	if !termutil.Isatty(os.Stdin.Fd()) {
		valByteArr, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		val = string(valByteArr)
	}

	return val
}
