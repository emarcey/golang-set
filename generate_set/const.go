package main

import (
	"reflect"
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

	KIND_TO_SET_TYPE_NAME = map[reflect.Kind]string{
		reflect.Array: ARRAY_SET_TYPE_NAME,
		reflect.Chan:  CHAN_SET_TYPE_NAME,
		reflect.Ptr:   PTR_SET_TYPE_NAME,
		reflect.Slice: SLICE_SET_TYPE_NAME,
		reflect.Map:   MAP_SET_TYPE_NAME,
	}
)
