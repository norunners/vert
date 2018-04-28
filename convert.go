// Package vert generically converts js objects to go values of any type with reflection.
// Additionally, go functions can be called with js object arguments.
package vert

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"reflect"
)

type Converter struct {
	typ reflect.Type
}

// New creates a new converter which generically converts js objects to go values of any type with reflection.
func New(value interface{}) *Converter {
	typ := typeOf(value)
	return &Converter{typ: typ}
}

// Interface converts the given source js object to an empty interface
// which can be type asserted back to the original type.
// Note, struct types must embed a js object pointer or be of type js object to allow for conversion.
// Beneficially, struct types can contain any depth of complex struct types.
// Note, map keys must be of type string to allow for conversion.
func (converter *Converter) Interface(src *js.Object) (interface{}, error) {
	value, err := convert(converter.typ, src)
	if err != nil {
		return nil, err
	}
	return value.Interface(), nil
}

// Value converts the given source js object to a reflection value.
func (converter *Converter) Value(src *js.Object) (reflect.Value, error) {
	return convert(converter.typ, src)
}

// typeOf returns the reflection type of the given value.
// Reflection type values are allowed.
func typeOf(value interface{}) reflect.Type {
	if typ, ok := value.(reflect.Type); ok {
		return typ
	}
	return reflect.TypeOf(value)
}

// convert recursively handles conversion of any type.
// Note, reflection panics are recovered and captured in an error.
func convert(typ reflect.Type, src *js.Object) (dst reflect.Value, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
	}()

	convert := option(typ)
	dst, err = convert(typ, src)
	return dst, err
}

// option returns the related convert function from the kind of the given type.
func option(typ reflect.Type) func(typ reflect.Type, src *js.Object) (reflect.Value, error) {
	switch typ.Kind() {
	case reflect.Ptr:
		return convertPointer
	case reflect.Struct:
		return convertStruct
	case reflect.Map:
		return convertMap
	case reflect.Slice:
		return convertSlice
	default:
		return convertDefault
	}
}

// convertPointer handles pointer types.
func convertPointer(typ reflect.Type, src *js.Object) (reflect.Value, error) {
	elemType := typ.Elem()
	dst := reflect.New(elemType)
	value, err := convert(elemType, src)
	if err != nil {
		return reflect.Zero(typ), err
	}
	dst.Elem().Set(value)
	return dst, nil
}

// convertStruct handles struct types.
// Note, struct types must embed a js object pointer or be of type js object.
func convertStruct(typ reflect.Type, src *js.Object) (reflect.Value, error) {
	dst := reflect.New(typ).Elem()
	if _, ok := dst.Interface().(js.Object); ok {
		objVal := reflect.ValueOf(src).Elem()
		dst.Set(objVal)
	} else {
		embedded := dst.Field(0)
		objVal, err := convert(embedded.Type(), src)
		if err != nil {
			return reflect.Zero(typ), err
		}
		embedded.Set(objVal)
	}
	return dst, nil
}

// convertMap handles map types.
// Note, map keys must be of type string.
func convertMap(typ reflect.Type, src *js.Object) (reflect.Value, error) {
	if srcType := reflect.TypeOf(src.Interface()); srcType.Kind() != reflect.Map {
		return reflect.Zero(typ), &typeErr{expected: reflect.Map, actual: src.Interface()}
	}

	valType := typ.Elem()
	keys := js.Keys(src)
	dst := reflect.MakeMapWithSize(typ, len(keys))
	for _, key := range keys {
		object := src.Get(key)
		value, err := convert(valType, object)
		if err != nil {
			return reflect.Zero(typ), err
		}
		keyVal := reflect.ValueOf(key)
		dst.SetMapIndex(keyVal, value)
	}
	return dst, nil
}

// convertSlice handles slice types.
func convertSlice(typ reflect.Type, src *js.Object) (reflect.Value, error) {
	srcType := reflect.TypeOf(src.Interface())
	if srcType.Kind() != reflect.Slice {
		return reflect.Zero(typ), &typeErr{expected: reflect.Slice, actual: src.Interface()}
	}

	elemType := typ.Elem()
	n := src.Length()
	dst := reflect.MakeSlice(typ, 0, n)
	for i := 0; i < n; i++ {
		object := src.Index(i)
		elem, err := convert(elemType, object)
		if err != nil {
			return reflect.Zero(typ), err
		}
		dst = reflect.Append(dst, elem)
	}
	return dst, nil
}

// convertDefault handles all other types, e.g. primitives.
// Note, convertDefault panics on incompatible type conversions.
func convertDefault(typ reflect.Type, src *js.Object) (reflect.Value, error) {
	srcVal := reflect.ValueOf(src.Interface())
	dst := srcVal.Convert(typ)
	return dst, nil
}
