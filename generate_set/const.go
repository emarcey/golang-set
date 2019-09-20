package main

import (
	"regexp"
)

const (
	STARTS_WITH_NUM_REGEX   = `^[0-9].*`
	SPLIT_OBJECT_NAME_REGEX = `[:/.\-_]+`

	ARRAY_SET_TYPE_NAME = "ArrayOf%v"
	CHAN_SET_TYPE_NAME  = "ChannelOf%v"
	PTR_SET_TYPE_NAME   = "PointerOf%v"
	SLICE_SET_TYPE_NAME = "SliceOf%v"
	MAP_SET_TYPE_NAME   = "MapOf%vTo%v"

	BASE_FILEPATH = "sets/%v_set"

	ITERATOR_FILENAME     = "%v_iterator.go"
	SET_FILENAME          = "%v_set.go"
	THREADSAFE_FILENAME   = "%v_threadsafe.go"
	THREADUNSAFE_FILENAME = "%v_threadunsafe.go"

	ITERATOR_TEMPLATE     = "generate_set/templates/iterator.gotemplate"
	SET_TEMPLATE          = "generate_set/templates/set.gotemplate"
	THREADSAFE_TEMPLATE   = "generate_set/templates/threadsafe.gotemplate"
	THREADUNSAFE_TEMPLATE = "generate_set/templates/threadunsafe.gotemplate"
)

var (
	STARTS_WITH_NUM_REGEXP   = regexp.MustCompile(STARTS_WITH_NUM_REGEX)
	SPLIT_OBJECT_NAME_REGEXP = regexp.MustCompile(SPLIT_OBJECT_NAME_REGEX)

	DEFAULT_TYPES = map[string]interface{}{
		"bool":    false,
		"int":     0,
		"int8":    0,
		"int16":   0,
		"int32":   0,
		"int64":   0,
		"uint":    0,
		"uint8":   0,
		"uint16":  0,
		"uint32":  0,
		"uint64":  0,
		"float32": 0.0,
		"float64": 0.0,
		"string":  `""`,
	}
)
