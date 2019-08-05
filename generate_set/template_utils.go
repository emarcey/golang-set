package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
)

var FUNC_MAP = template.FuncMap{
	"ToLower": strings.ToLower,
}

func CreateTemplate(field interface{}, templateBytes []byte, targetFileName string) error {
	genCode, err := ApplyTemplate(string(templateBytes), targetFileName, field, FUNC_MAP)
	if err != nil {
		return fmt.Errorf("cannot render template: %s: %v", templateBytes, err)
	}
	codeBytes, err := ioutil.ReadAll(genCode)
	if err != nil {
		return fmt.Errorf("cannot render template: %s: %v", templateBytes, err)
	}

	formattedCode, err := FormatCode(codeBytes)
	if err != nil {
		return fmt.Errorf("cannot render template: %s: %v", templateBytes, err)
	}
	return WriteGenFile(bytes.NewReader(formattedCode), targetFileName)
}

// ApplyTemplate is a helper methods that packages can call to render a
// template with any data and func map
func ApplyTemplate(templ string, templName string, data interface{}, funcMap template.FuncMap) (io.Reader, error) {
	codeTemplate, err := template.New(templName).Funcs(funcMap).Parse(templ)
	if err != nil {
		return nil, fmt.Errorf("cannot create template: %v", err)
	}

	outputBuffer := bytes.NewBuffer(nil)
	err = codeTemplate.Execute(outputBuffer, data)
	if err != nil {
		return nil, fmt.Errorf("template error: %v", err)
	}

	return outputBuffer, nil
}

func FormatCode(code []byte) ([]byte, error) {
	formatted, err := format.Source(code)
	if err != nil {
		return nil, fmt.Errorf("Code formatting error: %v", err)
	}
	return formatted, nil
}

func CreateSetFileFromTemplate(setType SetType, templateType TemplateType) error {
	template, err := GetTemplate(templateType.TemplateFilename)
	if err != nil {
		return err
	}
	targetFilename, err := makeFilename(setType.TitleName, templateType.OutFilename)
	if err != nil {
		return err
	}
	return CreateTemplate(setType, template, targetFilename)
}

func CreateSet(setType SetType, templateTypes []TemplateType) error {
	for _, templateType := range templateTypes {
		err := CreateSetFileFromTemplate(setType, templateType)
		if err != nil {
			return err
		}
	}
	return nil
}
