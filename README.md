Improved Go (igo)
===============

Everyone knows that Go is a very verbose language. It takes numerous lines of code to do what a few lines of code can do in other languages. This is a deliberate design decision by the Go Authors.

The igo project provides various syntactical sugar to make your code simpler and easier to read. It works by allowing you to program in `*.igo` files with the fancy new syntax. You then run `igo build` to transpile your `igo` files to standard `go` files which you can build as per normal.

1. Address Operator
    * Constants and Functions
2. Defers for for loops
    * `fordefer` guarantees to run prior to the loop's current iteration exiting.
3. Defer go
    * Run defer statements in a goroutine
  
This tool is an **experimental** is still at alpha stage.


## What is included

* igofmt
* igo compiler

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

Most professional front-end developers are fed up with standard JavaScript. They program using Typescript and then tranpile the code to standard ES5 JavaScript. igo adds the same step to the build process.

## Performance Tips

* Use `-m` command line flag to instruct gopher.js to minify code. Then minify further with [Webpack/UglifyJS](https://github.com/gopherjs/gopherjs/issues/136). A Webpack tutorial can be [found here](https://medium.com/ag-grid/webpack-tutorial-understanding-how-it-works-f73dfa164f01).
* Apply [gzip compression](https://en.wikipedia.org/wiki/HTTP_compression)
* Use int instead of (u)int8/16/32/64
* Use float64 instead of float32
* Try to avoid importing `fmt` package (either directly or indirectly).
* Use **react.JSFn()** and use native javascript functions as much as possible.
* https://github.com/gopherjs/gopherjs/wiki/JavaScript-Tips-and-Gotchas
* See if [jsgo](https://github.com/dave/jsgo) is appropriate for your web-based project.

## Future Work

* Fork [mapstructure](https://github.com/mitchellh/mapstructure) and remove dependecy for `fmt` and `net` (short term)
* Remove dependency for [mapstructure](https://github.com/mitchellh/mapstructure) (long term)
* WebAssembly version

#

### Legal Information

The license is a modified MIT license. Refer to `LICENSE` file for more details.

**© 2018 PJ Engineering and Business Solutions Pty. Ltd.**

### Final Notes

Feel free to enhance features by issuing pull-requests.

**Star** the project to show your appreciation.
