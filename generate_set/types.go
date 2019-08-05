package main

import (
	"fmt"
	"strings"
)

type SetType struct {
	DataType     string
	TitleName    string
	ImportPath   string
	DefaultValue string
}

func (s1 SetType) Equal(s2 SetType) bool {
	return s1.DataType == s2.DataType && s1.TitleName == s2.TitleName && s1.ImportPath == s2.ImportPath
}

func NewSetType(dataType, importPath, defaultValue string) SetType {
	titleName := makeSetTypeTitleName(dataType)
	return SetType{
		DataType:     dataType,
		TitleName:    titleName,
		ImportPath:   importPath,
		DefaultValue: defaultValue,
	}
}

func makeSetTypeTitleName(dataTypeName string) string {
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

type TemplateType struct {
	TemplateFilename string
	OutFilename      string
}

func NewTemplateType(templateFilename, outFilename string) TemplateType {
	return TemplateType{
		TemplateFilename: templateFilename,
		OutFilename:      outFilename,
	}
}

func MakeTemplateTypes() []TemplateType {
	return []TemplateType{
		NewTemplateType(ITERATOR_TEMPLATE, ITERATOR_FILENAME),
		NewTemplateType(SET_TEMPLATE, SET_FILENAME),
		NewTemplateType(THREADSAFE_TEMPLATE, THREADSAFE_FILENAME),
		NewTemplateType(THREADUNSAFE_TEMPLATE, THREADUNSAFE_FILENAME),
	}
}
