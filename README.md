vert
====
[![GoDoc](https://godoc.org/github.com/norunners/vert?status.svg)](https://godoc.org/github.com/norunners/vert)
[![Build Status](https://travis-ci.org/norunners/vert.svg?branch=master)](https://travis-ci.org/norunners/vert)

Package `vert` generically convert js objects to go values of any type and calls go functions with js object arguments.

Install
-------
```bash
go get github.com/norunners/vert
```

Converter
---------
Converter converts js objects to go values of any type.
##### Trivial hello world example:
The js object is converted to a string value and printed.
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
The js object is converted to the struct value `Data` and the `Message` field is printed.
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
	data := &Data{Object: js.Global.Get("Object").New()}
	data.Message = "Hello World!"
	return data.Object
}
```

Function Calls
--------------
Func calls go functions with js object arguments.
##### Trivial hello world example:
The vert function calls `print` with the string value converted from the js object.
```go
func main() {
    object := js.Global.Get("Object").New("Hello World!")
    function := vert.NewFunc(print)
    function.Call(object)
}

func print(message string) {
    fmt.Println(message)
}
```

##### Hello world example using struct values:
The vert function calls `print` with the struct value `Data` converted from the js object. 
```go
type Data struct {
    *js.Object
    Message string `js:"Message"`
}

func main() {
	object := jsData()
	function := vert.NewFunc(print)
	function.Call(object)
}

func print(data *Data) {
    fmt.Println(data.Message)
}

func jsData() *js.Object {
	data := &Data{Object: js.Global.Get("Object").New()}
	data.Message = "Hello World!"
	return data.Object
}
```

##### Addition example adds two numbers and returns the sum:
This vert function calls `add` with two int arguments and the result is type asserted to an int.
```go
func main() {
    a, b := 39, 3
    objA := js.Global.Get("Object").New(a)
    objB := js.Global.Get("Object").New(b)

    function := vert.NewFunc(add)
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
