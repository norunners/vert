//go:build js && wasm
// +build js,wasm

package vert

import (
	"reflect"
	"syscall/js"
	"testing"
	"unsafe"
)

func TestAssignToAll(t *testing.T) {
	v := ValueOf(allValue)
	ac := All{}
	err := v.AssignTo(&ac)
	if err != nil {
		t.Errorf("unexpected error: %+v\n", err)
	}
	if !reflect.DeepEqual(allValue, ac) {
		t.Errorf("expected: %+v but found: %+v\n", allValue, ac)
	}
}

func TestAssignToAllPointers(t *testing.T) {
	v := ValueOf(pallValue)
	ac := All{}
	err := v.AssignTo(&ac)
	if err != nil {
		t.Errorf("unexpected error: %+v\n", err)
	}
	if !reflect.DeepEqual(pallValue, ac) {
		t.Errorf("expected: %+v but found: %+v\n", pallValue, ac)
	}
}

func TestAssignToAllFields(t *testing.T) {
	exAll := reflect.ValueOf(allValue)
	acAll := reflect.ValueOf(&All{}).Elem()
	n := exAll.NumField()
	for i := 0; i < n; i++ {
		typ := exAll.Type()
		if f := exAll.Field(i); f.CanInterface() {
			name := typ.Field(i).Name
			t.Run(name, func(t *testing.T) {
				ex := f.Interface()
				ac := acAll.Field(i).Interface()

				v := ValueOf(ex)
				err := v.AssignTo(&ac)
				if err != nil {
					t.Errorf("unexpected error: %+v\n", err)
				}

				if !reflect.DeepEqual(ex, ac) {
					t.Errorf("expected: %+v but found: %+v\n", ex, ac)
				}
			})
		}
	}
}

func TestAssignToAllPointerFields(t *testing.T) {
	exAll := reflect.ValueOf(pallValue)
	acAll := reflect.ValueOf(&All{}).Elem()
	n := exAll.NumField()
	for i := 0; i < n; i++ {
		typ := exAll.Type()
		if f := exAll.Field(i); f.CanInterface() {
			name := typ.Field(i).Name
			t.Run(name, func(t *testing.T) {
				ex := f.Interface()
				ac := acAll.Field(i).Interface()

				v := ValueOf(ex)
				err := v.AssignTo(&ac)
				if err != nil {
					t.Errorf("unexpected error: %+v\n", err)
				}

				if !reflect.DeepEqual(ex, ac) {
					t.Errorf("expected: %+v but found: %+v\n", ex, ac)
				}
			})
		}
	}
}

type All struct {
	Bool    bool
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Uintptr uintptr
	Float32 float32
	Float64 float64

	Foo  string `js:"bar"`
	Foo2 string `js:"bar2"`

	PBool    *bool
	PInt     *int
	PInt8    *int8
	PInt16   *int16
	PInt32   *int32
	PInt64   *int64
	PUint    *uint
	PUint8   *uint8
	PUint16  *uint16
	PUint32  *uint32
	PUint64  *uint64
	PUintptr *uintptr
	PFloat32 *float32
	PFloat64 *float64

	String  string
	PString *string

	Map   map[string]Small
	MapP  map[string]*Small
	PMap  *map[string]Small
	PMapP *map[string]*Small

	EmptyMap map[string]Small
	NilMap   map[string]Small

	Slice   []Small
	SliceP  []*Small
	PSlice  *[]Small
	PSliceP *[]*Small

	EmptySlice []Small
	NilSlice   []Small

	StringSlice []string
	ByteSlice   []byte

	Small   Small
	PSmall  *Small
	PPSmall **Small

	Interface  interface{}
	PInterface *interface{}

	unexported int
}

type Small struct {
	Tag   string `js:"tag"`
	Field string
}

var allValue = All{
	Bool:    true,
	Int:     2,
	Int8:    3,
	Int16:   4,
	Int32:   5,
	Int64:   6,
	Uint:    7,
	Uint8:   8,
	Uint16:  9,
	Uint32:  10,
	Uint64:  11,
	Uintptr: 12,
	Float32: 14.1,
	Float64: 15.1,
	Foo:     "foo",
	Foo2:    "foo2",

	String: "16",
	Map: map[string]Small{
		"17": {Tag: "tag17", Field: "field17"},
		"18": {Tag: "tag18", Field: "field18"},
	},
	MapP: map[string]*Small{
		"19": {Tag: "tag19", Field: "field19"},
	},
	EmptyMap:    map[string]Small{},
	Slice:       []Small{{Tag: "tag20", Field: "field20"}, {Tag: "tag21", Field: "field21"}},
	SliceP:      []*Small{{Tag: "tag22", Field: "field22"}, {Tag: "tag23", Field: "field23"}},
	EmptySlice:  []Small{},
	StringSlice: []string{"str24", "str25", "str26"},
	ByteSlice:   []byte{27, 28, 29},
	Small:       Small{Tag: "tag30", Field: "field30"},
	PSmall:      &Small{Tag: "tag31", Field: "field31"},

	Interface: 5.2,
}

var pallValue = All{
	PBool:      &allValue.Bool,
	PInt:       &allValue.Int,
	PInt8:      &allValue.Int8,
	PInt16:     &allValue.Int16,
	PInt32:     &allValue.Int32,
	PInt64:     &allValue.Int64,
	PUint:      &allValue.Uint,
	PUint8:     &allValue.Uint8,
	PUint16:    &allValue.Uint16,
	PUint32:    &allValue.Uint32,
	PUint64:    &allValue.Uint64,
	PUintptr:   &allValue.Uintptr,
	PFloat32:   &allValue.Float32,
	PFloat64:   &allValue.Float64,
	PString:    &allValue.String,
	PMap:       &allValue.Map,
	PMapP:      &allValue.MapP,
	PSlice:     &allValue.Slice,
	PSliceP:    &allValue.SliceP,
	PPSmall:    &allValue.PSmall,
	PInterface: &allValue.Interface,
}

func TestInvalidAssignmentJsToGoError(t *testing.T) {
	var i interface{}
	tests := []struct {
		name string
		v    interface{}
		i    interface{}
		err  error
	}{
		{"bool to int", true, 0, &InvalidAssignmentError{Type: js.TypeBoolean, Kind: reflect.Int}},
		{"bool to float", true, 0.0, &InvalidAssignmentError{Type: js.TypeBoolean, Kind: reflect.Float64}},
		{"bool to string", true, "", &InvalidAssignmentError{Type: js.TypeBoolean, Kind: reflect.String}},

		{"int to bool", 1, false, &InvalidAssignmentError{Type: js.TypeNumber, Kind: reflect.Bool}},
		{"int to string", 2, "", &InvalidAssignmentError{Type: js.TypeNumber, Kind: reflect.String}},

		{"float to bool", 3.0, false, &InvalidAssignmentError{Type: js.TypeNumber, Kind: reflect.Bool}},
		{"float to string", 4.0, "", &InvalidAssignmentError{Type: js.TypeNumber, Kind: reflect.String}},

		{"string to bool", "str5", false, &InvalidAssignmentError{Type: js.TypeString, Kind: reflect.Bool}},
		{"string to int", "str6", 0, &InvalidAssignmentError{Type: js.TypeString, Kind: reflect.Int}},
		{"string to float", "str7", 0.0, &InvalidAssignmentError{Type: js.TypeString, Kind: reflect.Float64}},

		{"object to bool", object.New(), false, &InvalidAssignmentError{Type: js.TypeObject, Kind: reflect.Bool}},
		{"object to int", object.New(), 0, &InvalidAssignmentError{Type: js.TypeObject, Kind: reflect.Int}},
		{"object to float", object.New(), 0.0, &InvalidAssignmentError{Type: js.TypeObject, Kind: reflect.Float64}},
		{"object to complex", object.New(), 0i, &InvalidAssignmentError{Type: js.TypeObject, Kind: reflect.Complex128}},
		{"object to array", object.New(), [0]struct{}{}, &InvalidAssignmentError{Type: js.TypeObject, Kind: reflect.Array}},
		{"object to channel", object.New(), make(chan struct{}), &InvalidAssignmentError{Type: js.TypeObject, Kind: reflect.Chan}},
		{"object to func", object.New(), func() {}, &InvalidAssignmentError{Type: js.TypeObject, Kind: reflect.Func}},
		{"object to nil interface", object.New(), i, &InvalidAssignmentError{Type: js.TypeObject, Kind: reflect.Interface}},
		{"object to unsafe pointer", object.New(), unsafe.Pointer(uintptr(0)), &InvalidAssignmentError{Type: js.TypeObject, Kind: reflect.UnsafePointer}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := ValueOf(test.v)

			err := v.AssignTo(&test.i)

			if !reflect.DeepEqual(test.err, err) {
				t.Errorf("expected: %+v but found: %+v\n", test.err, err)
			}
		})
	}
}

func TestInvalidAssignmentNonNilPointerError(t *testing.T) {
	var p *struct{}
	tests := []struct {
		name string
		i    interface{}
		err  error
	}{
		{"nil pointer", p, &InvalidAssignmentError{Kind: reflect.Ptr}},
		{"bool", false, &InvalidAssignmentError{Kind: reflect.Bool}},
		{"int", 0, &InvalidAssignmentError{Kind: reflect.Int}},
		{"float", 0.0, &InvalidAssignmentError{Kind: reflect.Float64}},
		{"string", "", &InvalidAssignmentError{Kind: reflect.String}},
		{"struct", All{}, &InvalidAssignmentError{Kind: reflect.Struct}},
		{"map", map[struct{}]struct{}{}, &InvalidAssignmentError{Kind: reflect.Map}},
		{"slice", []struct{}{}, &InvalidAssignmentError{Kind: reflect.Slice}},
		{"invalid", nil, &InvalidAssignmentError{Kind: reflect.Invalid}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := &Value{}

			err := v.AssignTo(test.i)

			if !reflect.DeepEqual(test.err, err) {
				t.Errorf("expected: %+v but found: %+v\n", test.err, err)
			}
		})
	}
}
