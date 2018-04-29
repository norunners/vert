vert
====
[![GoDoc](https://godoc.org/github.com/norunners/vert?status.svg)](https://godoc.org/github.com/norunners/vert)
[![Build Status](https://travis-ci.org/norunners/vert.svg?branch=master)](https://travis-ci.org/norunners/vert)

Package `vert` generically converts js objects to go values of any type and calls go functions with js object arguments.

Install
-------
```bash
go get github.com/norunners/vert
```

Converter
---------
Converter generically converts js objects to go values of any type.
##### Trivial hello world example:
The js object is converted to the string value, then printed.
```go
func main() {
    object := js.Global.Get("Object").New("Hello World!")
    var message string
    converter := vert.New(message)
    result, _ := converter.Interface(object)
    message = result.(string)
    fmt.Println(message)
}
```

##### Hello world example using struct values:
The js object is converted to the `Data` struct value, then the `Message` field is printed.
```go
type Data struct {
    *js.Object
    Message string `js:"Message"`
}

func main() {
	object := jsData()
	var data *Data
	converter := vert.New(data)
	result, _ := converter.Interface(object)
	data = result.(*Data)
	fmt.Println(data.Message)
}

func jsData() *js.Object {
	object := js.Global.Get("Object").New()
	object.Set("Message", "Hello World!")
	return object
}
```

Function Calls
--------------
Converter calls go functions with js object arguments.
##### Trivial hello world example:
The converter calls the `print` function with the string value, which was converted from the js object.
```go
func main() {
    object := js.Global.Get("Object").New("Hello World!")
    function := vert.New(print)
    function.Call(object)
}

func print(message string) {
    fmt.Println(message)
}
```

##### Hello world example using struct values:
The converter calls the `print` function with the `Data` struct value, which was converted from the js object. 
```go
type Data struct {
    *js.Object
    Message string `js:"Message"`
}

func main() {
	object := jsData()
	function := vert.New(print)
	function.Call(object)
}

func print(data *Data) {
    fmt.Println(data.Message)
}

func jsData() *js.Object {
	object := js.Global.Get("Object").New()
	object.Set("Message", "Hello World!")
	return object
}
```

##### Addition example adds two numbers and returns the sum:
The converter calls the `add` function with two int arguments and the result is type asserted to an int.
```go
func main() {
    a, b := 39, 3
    objA := js.Global.Get("Object").New(a)
    objB := js.Global.Get("Object").New(b)

    function := vert.New(add)
    result, _ := function.Call(objA, objB)
    c := result[0].(int)

    fmt.Printf("%d + %d = %d\n", a, b, c)
}

func add(a, b int) int {
    return a + b
}
```

License
-------
* [MIT License](LICENSE)
