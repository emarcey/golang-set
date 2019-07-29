package main

import (
	"fmt"
	"reflect"
	"strings"
)

type SetType struct {
	DataType  string
	TitleName string
}

func NewSetType(dataType reflect.Type) SetType {
	return SetType{
		DataType:  dataType.Name(),
		TitleName: MakeSetTypeTitleName(dataType.Name()),
	}
}

func MakeSetTypeTitleName(dataTypeName string) string {
	if STARTS_WITH_NUM_REGEXP.MatchString(dataTypeName) {
		dataTypeName = fmt.Sprintf("_%v", dataTypeName)
	}
	return strings.ToTitle(dataTypeName)
}
