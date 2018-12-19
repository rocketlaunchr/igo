Improved Go (igo)
===============

Everyone knows that Go is a very verbose language. It takes numerous lines of code to do what a few lines of code can do in other languages. This is a deliberate design decision by the Go Authors.

The igo project provides various syntactical sugar to make your code simpler and easier to read. It works by allowing you to program in `*.igo` files with the fancy new syntax. You then run `igo build` to transpile your `igo` files to standard `go` files which you can then build as per normal.

1. Address Operator
    * Constants and Functions
2. Defers for for loops
    * `fordefer` guarantees to run prior to the loop's current iteration exiting.
3. Defer go
    * Run defer statements in a goroutine
  
This tool is **experimental** and still at alpha stage.


## What is included

* igofmt (auto format code)
* igo transpiler (generate standard go code)

## Installation

**Transpiler**

```
go get -u github.com/rocketlaunchr/igo
```

**Formatter**

```
go get -u github.com/rocketlaunchr/igo/igofmt
```

## Inspiration

Most professional front-end developers are fed up with standard JavaScript. They program using Typescript and then transpile the code to standard ES5 JavaScript. igo adds the same step to the build process.

## Examples

### Address Operator

```go

func main() {

	message := &"igo is so convenient"
	display(message)
   
	display(&`inline string`)

	def := &defaultMessage()
	display(def)

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

See https://blog.learngoprogramming.com/gotchas-of-defer-in-go-1-8d070894cb01 on why this is an improvement.

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

```go

mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	// Transmit how long request took to serve without delaying response to client.
	defer go transmitRequestStats(start)

	fmt.Fprintf(w, "Welcome to the home page!")
})

```

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



#

### Legal Information

The license is a modified MIT license. Refer to `LICENSE` file for more details.

**© 2018 PJ Engineering and Business Solutions Pty. Ltd.**

### Final Notes

Feel free to enhance features by issuing pull-requests.

**Star** the project to show your appreciation.
