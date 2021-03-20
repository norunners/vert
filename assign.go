// +build js,wasm

package vert

import (
	"fmt"
	"reflect"
	"syscall/js"
)

var zero = reflect.ValueOf(nil)

// AssignTo assigns a JS value to a Go pointer.
// Returns an error on invalid assignments.
func (v Value) AssignTo(i interface{}) error {
	rv := reflect.ValueOf(i)
	if k := rv.Kind(); k != reflect.Ptr || rv.IsNil() {
		return &InvalidAssignmentError{Kind: k}
	}

	return recoverAssignTo(rv, v.JSValue())
}

// recoverAssignTo recovers unexpected assignment panics.
// Please report unexpected panics.
func recoverAssignTo(rv reflect.Value, jv js.Value) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = &InvalidAssignmentError{rec: rec}
		}
	}()

	_, err = assignTo(rv, jv)
	return
}

// assignTo recursively assigns a value.
func assignTo(rv reflect.Value, jv js.Value) (reflect.Value, error) {
	if jv == js.Null() || jv == js.Undefined() {
		return zero, nil
	}

	k := rv.Kind()
	switch k {
	case reflect.Ptr:
		return assignToPointer(rv, jv)
	case reflect.Interface:
		if e := rv.Elem(); e != zero {
			return assignToInterface(rv, e, jv)
		}
	}

	switch t := jv.Type(); t {
	case js.TypeBoolean:
		return assignToBasic(rv, jv.Bool(), t)
	case js.TypeNumber:
		return assignToBasic(rv, jv.Float(), t)
	case js.TypeString:
		return assignToBasic(rv, jv.String(), t)
	case js.TypeObject:
		return assignToValue(rv, jv)
	default:
		return zero, &InvalidAssignmentError{Type: t, Kind: k}
	}
}

// assignToPointer assigns a value to a pointer.
func assignToPointer(p reflect.Value, jv js.Value) (reflect.Value, error) {
	if p.IsNil() {
		p = reflect.New(p.Type().Elem())
	}

	v, err := assignTo(p.Elem(), jv)
	if err != nil {
		return zero, err
	}
	if v != zero {
		p.Elem().Set(v)
	}
	return p, nil
}

// assignToInterface assigns a value to an interface.
func assignToInterface(i, e reflect.Value, jv js.Value) (reflect.Value, error) {
	v, err := assignTo(e, jv)
	if err != nil {
		return zero, err
	}
	if v != zero {
		i.Set(v)
	}
	return i, nil
}

// assignToBasic assigns a primitive value to a basic value.
func assignToBasic(b reflect.Value, i interface{}, t js.Type) (val reflect.Value, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = &InvalidAssignmentError{Type: t, Kind: b.Kind()}
		}
	}()

	v := reflect.ValueOf(i)
	val = v.Convert(b.Type())
	return
}

// assignToObject assigns an object to a value.
func assignToValue(rv reflect.Value, jv js.Value) (reflect.Value, error) {
	switch k := rv.Kind(); k {
	case reflect.Struct:
		return assignToStruct(rv, jv)
	case reflect.Map:
		return assignToMap(rv, jv)
	case reflect.Slice:
		return assignToSlice(rv, jv)
	case reflect.Array:
		return assignToArray(rv, jv)
	default:
		return zero, &InvalidAssignmentError{Type: jv.Type(), Kind: k}
	}
}

// assignToStruct assigns an object to a struct.
func assignToStruct(s reflect.Value, val js.Value) (reflect.Value, error) {
	t := s.Type()
	s = reflect.New(t).Elem()
	n := s.NumField()
	for i := 0; i < n; i++ {
		if f := s.Field(i); f.CanInterface() {
			k := nameOf(t.Field(i))
			jf := val.Get(k)
			v, err := assignTo(f, jf)
			if err != nil {
				return zero, err
			}
			if v == zero {
				continue
			}
			f.Set(v)
		}
	}
	return s, nil
}

// assignToMap assigns an object to a map.
// Map keys must be of type string.
func assignToMap(m reflect.Value, jv js.Value) (reflect.Value, error) {
	t := m.Type()
	keys := object.Call("keys", jv)
	n := keys.Length()
	if m.IsNil() {
		m = reflect.MakeMapWithSize(t, n)
	}
	kt := t.Key()
	vt := t.Elem()
	for i := 0; i < n; i++ {
		jk := keys.Index(i)
		k := reflect.New(kt).Elem()
		k, err := assignTo(k, jk)
		if err != nil {
			return zero, err
		}
		if k == zero {
			continue
		}
		jv := jv.Get(jk.String())
		v := reflect.New(vt).Elem()
		v, err = assignTo(v, jv)
		if err != nil {
			return zero, err
		}
		if v == zero {
			continue
		}
		m.SetMapIndex(k, v)
	}
	return m, nil
}

// assignToSlice assigns an array object to a slice.
func assignToSlice(s reflect.Value, jv js.Value) (reflect.Value, error) {
	t := s.Type()
	n := jv.Length()
	if s.IsNil() {
		s = reflect.MakeSlice(t, 0, n)
	}
	et := t.Elem()
	for i := 0; i < n; i++ {
		e := reflect.New(et).Elem()
		je := jv.Index(i)
		e, err := assignTo(e, je)
		if err != nil {
			return zero, err
		}
		if e == zero {
			continue
		}
		s = reflect.Append(s, e)
	}
	return s, nil
}

func assignToArray(s reflect.Value, jv js.Value) (reflect.Value, error) {
	t := s.Type()
	n := jv.Length()
	et := t.Elem()
	if !s.CanAddr() {
		s = reflect.New(reflect.ArrayOf(n, t.Elem())).Elem()
	}
	for i := 0; i < n; i++ {
		e := reflect.New(et).Elem()
		je := jv.Index(i)
		e, err := assignTo(e, je)
		if err != nil {
			return zero, err
		}
		if e == zero {
			continue
		}
		s.Index(i).Set(e)
	}
	return s, nil
}

type InvalidAssignmentError struct {
	Type js.Type
	Kind reflect.Kind
	rec  interface{}
}

func (e *InvalidAssignmentError) Error() string {
	if e.rec != nil {
		return fmt.Sprintf("unexpected panic: %+v", e.rec)
	}
	if e.Type == js.TypeUndefined {
		return fmt.Sprintf("invalid assignment to Go kind: %v must be a non-nil pointer", e.Kind)
	}
	return fmt.Sprintf("invalid assignment from JS type: %v to Go kind: %v", e.Type, e.Kind)
}
