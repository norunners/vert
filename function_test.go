package vert

import (
	"testing"
)

func TestCallFuncEmpty(t *testing.T) {
	function := NewFunc(testEmptyFunc)
	result, err := function.Call()

	assertNoError(t, err)
	assertEqual(t, 0, len(result))
}

func TestCallFuncEmptyWithArg(t *testing.T) {
	object := newObject("testString")

	function := NewFunc(testEmptyFunc)
	result, err := function.Call(object)

	assertNoError(t, err)
	assertEqual(t, 0, len(result))
}

func TestCallFuncString(t *testing.T) {
	expected := "testString"
	object := newObject(expected)

	function := NewFunc(testStringFunc)
	result, err := function.Call(object)

	assertNoError(t, err)
	actual := result[0].(string)
	assertEqual(t, expected, actual)
}

func TestCallFuncStruct(t *testing.T) {
	expected := &testData{Object: newObject()}
	expected.text = "testText"
	object := newObject(expected)

	function := NewFunc(testDataFunc)
	result, err := function.Call(object)

	assertNoError(t, err)
	actual := result[0].(*testData)
	assertEqual(t, expected.Object.Interface(), actual.Object.Interface())
	assertEqual(t, expected.text, actual.text)
}

func TestCallFuncNoArgs(t *testing.T) {
	function := NewFunc(testDataFunc)
	_, err := function.Call()

	assertError(t, err)
}

func TestCallFuncBadConvert(t *testing.T) {
	object := newObject(42)

	function := NewFunc(testStringFunc)
	_, err := function.Call(object)

	assertError(t, err)
}

func testEmptyFunc() {
}

func testStringFunc(value string) string {
	return value
}

func testDataFunc(value *testData) *testData {
	return value
}
