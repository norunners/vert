package vert

import (
	"fmt"
	"reflect"
)

type typeErr struct {
	expected reflect.Kind
	actual   interface{}
}

func (err *typeErr) Error() string {
	return fmt.Sprintf("expected: %s but found: %T in value: %v",
		err.expected.String(), err.actual, err.actual)
}
