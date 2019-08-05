package main

import (
	"flag"
	"fmt"
	"os"
)

func generateSet() error {
	var structName = flag.String("struct_name", "", "name of struct to generate")
	var importPath = flag.String("import_path", "", "go")
	var defaultValue = flag.String("default_value", "", "default value of struct")

	flag.Parse()
	if *structName == "" {
		return NewEmptyFlagError("struct_name")
	}
	if *defaultValue == "" {
		return NewEmptyFlagError("default_value")
	}

	setType := NewSetType(*structName, *importPath, *defaultValue)

	templateTypes := MakeTemplateTypes()

	return CreateSet(setType, templateTypes)
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
