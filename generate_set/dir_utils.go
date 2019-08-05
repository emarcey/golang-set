package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func makeFilename(titleName string, baseFilename string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	lowerTitleName := strings.ToLower(titleName)
	filePath := filepath.Join(wd, BASE_FILEPATH, baseFilename)
	return fmt.Sprintf(filePath, lowerTitleName, lowerTitleName), nil
}

func GetTemplate(templateFilename string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filename := filepath.Join(wd, templateFilename)
	return ioutil.ReadFile(filename)
}

func WriteGenFile(file io.Reader, path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(outFile, file)
	if err != nil {
		return err
	}
	return outFile.Close()
}
