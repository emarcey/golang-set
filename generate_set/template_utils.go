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

func (s1 SetType) Equal(s2 SetType) bool {
	return s1.DataType == s2.DataType && s1.TitleName == s2.TitleName
}

func NewSetType(dataType reflect.Type) (SetType, error) {
	titleName, err := makeSetTypeTitleName(dataType)
	if err != nil {
		return SetType{}, err
	}
	return SetType{
		DataType:  dataType.String(),
		TitleName: titleName,
	}, nil
}

func makeSetTypeTitleName(dataType reflect.Type) (string, error) {
	switch k := dataType.Kind(); k {
	case reflect.Array, reflect.Chan, reflect.Ptr, reflect.Slice:
		elemName, err := makeSetTypeTitleName(dataType.Elem())
		if err != nil {
			return "", err
		}
		setTypeName, ok := KIND_TO_SET_TYPE_NAME[k]
		if !ok {
			return "", fmt.Errorf("Unexpected Kind: %v\n", k)
		}
		return fmt.Sprintf(setTypeName, elemName), nil
	case reflect.Map:
		keyName, err := makeSetTypeTitleName(dataType.Key())
		if err != nil {
			return "", err
		}
		elemName, err := makeSetTypeTitleName(dataType.Elem())
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(MAP_SET_TYPE_NAME, keyName, elemName), nil
	default:
		return makeSimpleSetTypeTitleName(dataType.String()), nil
	}
}

func makeSimpleSetTypeTitleName(dataTypeName string) string {
	splitTitleName := SPLIT_OBJECT_NAME_REGEXP.Split(dataTypeName, -1)
	titleName := ""

	for i, s := range splitTitleName {
		if i == 0 && STARTS_WITH_NUM_REGEXP.MatchString(s) {
			s = fmt.Sprintf("_%v", s)
		}
		titleName += strings.Title(s)
	}
	return strings.TrimSpace(strings.TrimSuffix(titleName, "{}"))
}
