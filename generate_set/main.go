package main

import (
	"flag"
	"fmt"
	"os"
)

func generateSet() error {
	var structName = flag.String("struct_name", "", "name of struct to generate")
	var importPath = flag.String("import_path", "", "go")
	if *structName == "" {
		return NewEmptyFlagError("struct_name")
	}

	_ = NewSetType(*structName, *importPath)

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
