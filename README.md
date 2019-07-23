Improved Go (igo)
===============

Everyone knows that Go is a very verbose language. It takes numerous lines of code to do what a few lines of code can do in other languages. This is a deliberate design decision by the Go Authors.

The igo project provides various syntactical sugar to make your code simpler and easier to read. It works by allowing you to program in `*.igo` files with the fancy new syntax. You then run `igo build` to transpile your `igo` files to standard `go` files which you can then build as per normal.

1. Address Operator (&)
    * Constants and Functions
2. Defers for for-loops
    * `fordefer` guarantees to run prior to the loop's current iteration exiting.
3. Defer go
    * Run defer statements in a goroutine
4. must function
    * Converts a multi-return value function into a single-return function.
    * See [#32219](https://github.com/golang/go/issues/32219)


**NOTE: igo is pronounced ee-gohr**
  

## What is included

* igofmt (auto format code)
* igo transpiler (generate standard go code)

## Installation

**Transpiler**

```
go get -u github.com/rocketlaunchr/igo
```

Use `go install` to install the executable. 

**Formatter**

```
go get -u github.com/rocketlaunchr/igo/igofmt
```

## Inspiration

Most professional front-end developers are fed up with standard JavaScript. They program using Typescript and then transpile the code to standard ES5 JavaScript. igo adds the same step to the build process.

## Examples

### Address Operator

The Address Operator allows you to use more visually pleasing syntax. There is no need for a temporary variable. It can be used with `string`, `bool`, `int`, `float64` and function calls where the function returns 1 return value.


```go

func main() {

	message := &"igo is so convenient"
	display(message)
   
	display(&`inline string`)

	display(&defaultMessage())

}

func display(m *string) {
	if m == nil {
		fmt.Print("no message")
	} else {
		fmt.Print(*m)
	}

}

func defaultMessage() string {
	return "default message"
}

```

### Fordefer

See [Blog post](https://blog.learngoprogramming.com/gotchas-of-defer-in-go-1-8d070894cb01) on why this is an improvement. It can be especially helpful in unit tests.

```go

for {
	row, err := db.Query("SELECT ...")
	if err != nil {
		panic(err)
	}

	fordefer row.Close()
}

```


### Defer go

This feature makes Go's language syntax more internally consistent. There is no reason why `defer` and `go` should not work together.

```go

mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	// Transmit how long the request took to serve without delaying response to client.
	defer go transmitRequestStats(start)

	fmt.Fprintf(w, "Welcome to the home page!")
})

```

### Must builtin function

`must` is a "builtin" function that converts a multi-return value function (`"fn"`) into a single-return function. `fn's` final return value is expected to be of type `error`. `must` will panic upon encountering an error.

It is useful in scenarios where you know that no error will actually be returned by `fn` and you just want to use the function inline. Alternatively, you may want to catch the error during local development because no error should be produced in production.

`must` also accepts an optional second argument of type `func(error) error`.

See [#32219](https://github.com/golang/go/issues/32219)


```go
import "database/sql"

db := must(sql.Open("mysql", "host"))

```

**LIMITATIONS**

* Currently, it only works when `fn` returns **two** return values.
* It doesn't work when used outside of functions (i.e. initializing package variables).
* It works perfectly in simple cases. For more complex cases, peruse the generated code.
* A PR would be appreciated by an expert in the `go/types` package. It is possible to create a truly generics-compatible `must` that resolves the limitations above.
* Unlike real "builtin" functions, `must` is a reserved keyword.

## How to use

### Transpile

`igo` can accept numerous directories or igo files. The generated go files are saved alongside the igo files.

```
igo build [igo files...]
```

### Format Code

`igofmt` will format your code to the standard form. It understands igo syntax.

```
igofmt [-s] [igo files...]
```

Configure your IDE to run `igofmt` upon saving a `*.igo` file.

## Design Decisions and Limitations

Pull-Requests are requested for the below deficiencies.

* For `fordefer`: `goto` statements inside a for-loop that jump outside the for-loop is not implemented. Use `github.com/rocketlaunchr/igo/stack` package manually in such cases.
* `igofmt -s` Simplified mode is not implemented. [See here for instructions on issuing a PR](https://github.com/golang/go/blob/master/src/cmd/gofmt/simplify.go#L15).
* `goimports` equivalent has not been made.
* Address Operator for constants currently only supports `string`, `bool`, `float64` and `int`. The other int types are not supported. This can be fixed by using [go/types](https://github.com/golang/example/tree/master/gotypes) package.
* Address Operator feature assumes you have not attempted to redefine `true` and `false` to something/anything else.
	* Why would you redefine them anyway?

## Tips & Advice

* Store the `igo` and generated `go` files in your git repository.
* Configure your IDE to run `igofmt` upon saving a `*.igo` file.

#

### Legal Information

The license is a modified MIT license. Refer to the `LICENSE` file for more details.

**© 2018-19 PJ Engineering and Business Solutions Pty. Ltd.**

### Final Notes

Feel free to enhance features by issuing pull-requests.

**Star** the project to show your appreciation.
