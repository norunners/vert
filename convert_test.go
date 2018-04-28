package vert

import (
	"github.com/gopherjs/gopherjs/js"

	"reflect"
	"testing"
)

type testData struct {
	*js.Object

	text string `js:"text"`
}

type testComplex struct {
	*js.Object

	inner *testData           `js:"inner"`
	data  map[string]testData `js:"data"`
}

func TestConvertStruct(t *testing.T) {
	expected := &testData{Object: newObject()}
	expected.text = "testText"
	object := newObject(expected)

	var actual *testData
	converter := New(actual)
	result, err := converter.Interface(object)
	assertNoError(t, err)
	actual = result.(*testData)

	assertEqual(t, expected.Object.Interface(), actual.Object.Interface())
	assertEqual(t, expected.text, actual.text)
}

func TestConvertObject(t *testing.T) {
	expected := &testData{Object: newObject()}
	expected.text = "world"
	object := newObject(expected)

	var actual *js.Object
	converter := New(actual)
	result, err := converter.Interface(object)

	assertNoError(t, err)
	actual = result.(*js.Object)
	assertEqual(t, expected.Interface(), actual.Interface())
}

func TestConvertPointer(t *testing.T) {
	number := "testPointer"
	expected := &number
	object := newObject(expected)

	var actual *string
	converter := New(actual)
	result, err := converter.Interface(object)

	assertNoError(t, err)
	actual = result.(*string)
	assertEqual(t, expected, actual)
}

func TestConvertMap(t *testing.T) {
	expected := make(map[string]int)
	expected["answer"] = 42
	object := newObject(expected)

	var actual map[string]int
	converter := New(actual)
	result, err := converter.Interface(object)

	assertNoError(t, err)
	actual = result.(map[string]int)
	assertEqual(t, expected, actual)
}

func TestConvertSlice(t *testing.T) {
	expected := make([]bool, 0)
	expected = append(expected, false)
	expected = append(expected, true)
	object := newObject(expected)

	var actual []bool
	converter := New(actual)
	result, err := converter.Interface(object)

	assertNoError(t, err)
	actual = result.([]bool)
	assertEqual(t, expected, actual)
}

func TestConvertString(t *testing.T) {
	expected := "testString"
	object := newObject(expected)

	var actual string
	converter := New(actual)
	result, err := converter.Interface(object)

	assertNoError(t, err)
	actual = result.(string)
	assertEqual(t, expected, actual)
}

func TestConvertBool(t *testing.T) {
	expected := true
	object := newObject(expected)

	var actual bool
	converter := New(actual)
	result, err := converter.Interface(object)

	assertNoError(t, err)
	actual = result.(bool)
	assertEqual(t, expected, actual)
}

func TestConvertInt(t *testing.T) {
	expected := 42
	object := newObject(expected)

	var actual int
	converter := New(actual)
	result, err := converter.Interface(object)

	assertNoError(t, err)
	actual = result.(int)
	assertEqual(t, expected, actual)
}

func TestConvertFloat64(t *testing.T) {
	expected := 3.14
	object := newObject(expected)

	var actual float64
	converter := New(actual)
	result, err := converter.Interface(object)

	assertNoError(t, err)
	actual = result.(float64)
	assertEqual(t, expected, actual)
}

func TestConvertComplexStruct(t *testing.T) {
	inner := &testData{Object: newObject()}
	inner.text = "testText"
	data := &testData{Object: newObject()}
	data.text = "testKey"

	expected := &testComplex{Object: newObject()}
	expected.inner = inner
	expected.data = make(map[string]testData)
	expected.data["testKey"] = *data

	object := newObject(expected)

	var actual *testComplex
	converter := New(actual)
	result, err := converter.Interface(object)

	assertNoError(t, err)
	actual = result.(*testComplex)
	assertEqual(t, expected.Object.Interface(), actual.Object.Interface())
	assertEqual(t, expected.data, actual.data)
	assertEqual(t, expected.inner.Object.Interface(), actual.inner.Object.Interface())
	assertEqual(t, expected.inner.text, actual.inner.text)
}

func TestConvertComplexObject(t *testing.T) {
	inner := &testData{Object: newObject()}
	inner.text = "testText"
	data := &testData{Object: newObject()}
	data.text = "testValue"

	expected := &testComplex{Object: newObject()}
	expected.inner = inner
	expected.data = make(map[string]testData)
	expected.data["testKey"] = *data
	object := newObject(expected)

	var actual *js.Object
	converter := New(actual)
	result, err := converter.Interface(object)

	assertNoError(t, err)
	actual = result.(*js.Object)
	assertEqual(t, expected.Interface(), actual.Interface())
}

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("expected error but found nil")
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("expected nil but found error: %v", err)
	}
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %#v but found: %#v", expected, actual)
	}
}

func newObject(args ...interface{}) *js.Object {
	return js.Global.Get("Object").New(args...)
}
