[![Go Report Card](https://goreportcard.com/badge/github.com/emarcey/golang-set)](https://goreportcard.com/report/github.com/emarcey/golang-set)
[![GoDoc](https://godoc.org/github.com/emarcey/golang-set?status.svg)](http://godoc.org/github.com/emarcey/golang-set)

## golang-set

Fork of deckarep's [golang-set](https://github.com/deckarep/golang-set) with generate code to make type-specific Sets, because I was sick of converting in and out of interfaces.

Comes with a bunch of sets based on basic types.

Doesn't support Cartesian Products or Power Sets yet.

### Examples

To build
```
go build -o generate_set_exec ./generate_set
```

To generate a set for the `time.Time` object
```
./generate_set_exec -struct_name="time.Time" -import_path="time" -default_value="time.Time{}"
```

To generate a bunch of basic types
```
./generate_set_exec -make_defaults=`true`
```
