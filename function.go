package vert

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"reflect"
)

type Func struct {
	val reflect.Value
	typ reflect.Type
}

// NewFunc creates a new function which can be called with js object arguments.
// Note, functions with zero arguments are also allowed.
func NewFunc(function interface{}) *Func {
	val := reflect.ValueOf(function)
	typ := val.Type()
	return &Func{val: val, typ: typ}
}

// Call calls the function with the given js object arguments.
// First, the arguments are converted to their respective types from the function arguments.
// The return values are returned as a slice of empty interfaces
// which can be type asserted back to their original types.
// An error is returned if any argument fails to convert
// or there are not enough arguments to satisfy the number of function arguments.
// Note, the function is called with the arguments in the order given
// and it is safe to provide more arguments than the number of function arguments.
func (function *Func) Call(args ...*js.Object) ([]interface{}, error) {
	argVals, err := function.Values(args...)
	if err != nil {
		return nil, err
	}
	values := function.val.Call(argVals)
	interfaces := make([]interface{}, 0, len(values))
	for _, value := range values {
		interfaces = append(interfaces, value.Interface())
	}
	return interfaces, nil
}

// Values converts the given js object arguments to a slice of reflection values from the types of the function arguments.
// An error is returned if any argument fails to convert
// or there are not enough arguments to satisfy the number of function arguments.
func (function *Func) Values(args ...*js.Object) ([]reflect.Value, error) {
	n := function.typ.NumIn()
	length := len(args)
	if length < n {
		return nil, fmt.Errorf("expected at least: %d js objects but found: %d", n, length)
	}
	values := make([]reflect.Value, 0, n)
	for i := 0; i < n; i++ {
		argType := function.typ.In(i)
		value, err := convert(argType, args[i])
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}