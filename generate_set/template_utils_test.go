package main

import (
	"reflect"
	"testing"
	"time"
)

func TestMakeNewSetType(t *testing.T) {
	var testCases = []struct {
		given    reflect.Type
		expected SetType
	}{
		{
			given: reflect.TypeOf(int(1)),
			expected: SetType{
				DataType:  "int",
				TitleName: "Int",
			},
		},
		{
			given: reflect.TypeOf(int64(1)),
			expected: SetType{
				DataType:  "int64",
				TitleName: "Int64",
			},
		},
		{
			given: reflect.TypeOf(time.Time{}),
			expected: SetType{
				DataType:  "time.Time",
				TitleName: "TimeTime",
			},
		},
		{
			given: reflect.TypeOf([]interface{}{}),
			expected: SetType{
				DataType:  "[]interface {}",
				TitleName: "SliceOfInterface",
			},
		},
	}

	for i, testCase := range testCases {
		result, _ := NewSetType(testCase.given)
		if !testCase.expected.Equal(result) {
			t.Error("test", i, "given", testCase.given, "expected", testCase.expected, "result", result)
		}
	}
}

func TestMakeSimpleSetTypeTitleName(t *testing.T) {
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
		result := makeSimpleSetTypeTitleName(testCase.given)
		if testCase.expected != result {
			t.Error("test", i, "given", testCase.given, "expected", testCase.expected, "result", result)
		}
	}
}
