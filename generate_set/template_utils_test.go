package main

import (
	"testing"
)

func TestMakeNewSetType(t *testing.T) {
	var testCases = []struct {
		givenType       string
		givenImportPath string
		expected        SetType
	}{
		{
			givenType:       "int",
			givenImportPath: "",
			expected: SetType{
				DataType:  "int",
				TitleName: "Int",
			},
		},
		{
			givenType:       "int64",
			givenImportPath: "thing/thing/thing",
			expected: SetType{
				DataType:   "int64",
				TitleName:  "Int64",
				ImportPath: "thing/thing/thing",
			},
		},
		{
			givenType:       "time.Time",
			givenImportPath: "",
			expected: SetType{
				DataType:  "time.Time",
				TitleName: "TimeTime",
			},
		},
		{
			givenType:       "interface{}",
			givenImportPath: "",
			expected: SetType{
				DataType:  "interface{}",
				TitleName: "Interface",
			},
		},
	}

	for i, testCase := range testCases {
		result := NewSetType(testCase.givenType, testCase.givenImportPath)
		if !testCase.expected.Equal(result) {
			t.Error("test", i, "given", testCase.givenType, "and", testCase.givenImportPath, "expected", testCase.expected, "result", result)
		}
	}
}

func TestMakeSetTypeTitleName(t *testing.T) {
	var testCases = []struct {
		given    string
		expected string
	}{
		{
			given:    "",
			expected: "",
		},
		{
			given:    ".",
			expected: "",
		},
		{
			given:    "test64",
			expected: "Test64",
		},
		{
			given:    "testMeOut64",
			expected: "TestMeOut64",
		},
		{
			given:    "64test",
			expected: "_64test",
		},
		{
			given:    "test/me.oUt-this:time64",
			expected: "TestMeOUtThisTime64",
		},
	}

	for i, testCase := range testCases {
		result := makeSetTypeTitleName(testCase.given)
		if testCase.expected != result {
			t.Error("test", i, "given", testCase.given, "expected", testCase.expected, "result", result)
		}
	}
}
