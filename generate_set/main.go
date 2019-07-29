package main

import (
	"fmt"
	"os"
)

func generateSet() error {
	return nil
}

func main() {
	var exitcode int = 0

	err := generateSet()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error, Exiting: %v", err))
		exitcode = 1
	}
	os.Exit(exitcode)
}
