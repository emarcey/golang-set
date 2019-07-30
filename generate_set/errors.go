package main

import (
	"fmt"
	"reflect"
)

func NewEmptyFlagError(flagName string) error {
	return fmt.Errorf("Received empty value for flag %v.", flagName)
}

func NewTypeNotSupportedError(t reflect.Type) error {
	return fmt.Errorf("Type %v not supported.", t)
}
