// +build js,wasm

// Package vert provides WebAssembly interop between Go and JS values.
package vert

import (
	"reflect"
	"syscall/js"
)

var (
	null   = js.ValueOf(nil)
	object = js.Global().Get("Object")
	array  = js.Global().Get("Array")
)

// Value is an assignable JS value.
type Value struct {
	js.Value
}

// JSValue returns the JS value.
func (v Value) JSValue() js.Value {
	return v.Value
}

// ValueOf returns the Go value as a new value.
func ValueOf(i interface{}) Value {
	switch i.(type) {
	case nil, js.Value, js.Wrapper:
		return Value{Value: js.ValueOf(i)}
	default:
		v := reflect.ValueOf(i)
		return Value{Value: valueOf(v)}
	}
}

// valueOf recursively returns a new value.
func valueOf(v reflect.Value) js.Value {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return valueOfPointerOrInterface(v)
	case reflect.Slice:
		return valueOfSlice(v)
	case reflect.Array:
		return valueOfArray(v)
	case reflect.Map:
		return valueOfMap(v)
	case reflect.Struct:
		return valueOfStruct(v)
	default:
		return js.ValueOf(v.Interface())
	}
}

// valueOfPointerOrInterface returns a new value.
func valueOfPointerOrInterface(v reflect.Value) js.Value {
	if v.IsNil() {
		return null
	}
	return valueOf(v.Elem())
}

// valueOfSlice returns a new array object value.
func valueOfSlice(v reflect.Value) js.Value {
	if v.IsNil() {
		return null
	}
	return valueOfArray(v)
}

// valueOfArray returns a new array object value.
func valueOfArray(v reflect.Value) js.Value {
	a := array.New()
	n := v.Len()
	for i := 0; i < n; i++ {
		e := v.Index(i)
		a.SetIndex(i, valueOf(e))
	}
	return a
}

// valueOfMap returns a new object value.
// Map keys must be of type string.
func valueOfMap(v reflect.Value) js.Value {
	if v.IsNil() {
		return null
	}
	m := object.New()
	i := v.MapRange()
	for i.Next() {
		k := i.Key().Interface().(string)
		m.Set(k, valueOf(i.Value()))
	}
	return m
}

// valueOfStruct returns a new object value.
func valueOfStruct(v reflect.Value) js.Value {
	t := v.Type()
	s := object.New()
	n := v.NumField()
	for i := 0; i < n; i++ {
		if f := v.Field(i); f.CanInterface() {
			k := nameOf(t.Field(i))
			s.Set(k, valueOf(f))
		}
	}
	return s
}

// nameOf returns the JS tag name, otherwise the field name.
func nameOf(sf reflect.StructField) string {
	name := sf.Tag.Get("js")
	if name == "" {
		return sf.Name
	}
	return name
}
