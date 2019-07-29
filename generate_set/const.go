package main

import (
	"regexp"
)

const (
	STARTS_WITH_NUM_REGEX = `^[0-9].*`
)

var (
	STARTS_WITH_NUM_REGEXP = regexp.MustCompile(STARTS_WITH_NUM_REGEX)
)
