package util

import (
	"fmt"
	"reflect"
)

func isNil(value interface{}) bool {
	// https://mangatmodi.medium.com/go-check-nil-interface-the-right-way-d142776edef1
	if value == nil {
		return true
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(value).IsNil()
	}
	return false
}

func AssertNotNill(value interface{}, error interface{}) {
	if isNil(value) {
		print(value, value, value)
		panic(error)
	}
}

func AssertNill(value interface{}) {
	if !isNil(value) {
		panic(value)
	}
}

func AssertTrue(value bool, error interface{}) {
	if !value {
		panic(error)
	}
}

func AssertIntEqual(valueA int, valueB int, error string) {
	if valueA != valueB {
		panic(
			fmt.Sprintf("%s (%d != %d)", error, valueA, valueB),
		)
	}
}
