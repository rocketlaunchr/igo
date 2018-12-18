// Copyright 2018 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"go/format"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"

	"github.com/rocketlaunchr/igo/common"
)

// https://golang.org/cmd/gofmt/
// http://www.tothenew.com/blog/gofmt-formatting-the-go-code/
// https://spf13.com/post/go-fmt/

var (
	simplifyAST = flag.Bool("s", false, "simplify code (not implemented yet)")
	exitCode    = 0
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: igofmt [flags] [path ...]\n")
	flag.PrintDefaults()
}

func main() {
	gofmtMain()
	os.Exit(exitCode)
}

func gofmtMain() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		return
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)

		files, err := common.Files(path, "igo")
		if err != nil {
			exitCode = common.Report(err)
			return
		}

		for _, path := range files {
			err := processFile(path)
			if err != nil {
				exitCode = common.Report(err)
				return
			}
		}

	}
}

type state int

const (
	initial    state = 0
	foFound          = 1
	dotFound         = 2
	dfFound          = 3
	lBrakFound       = 4
	rBrakFound       = 5
)

func processFile(path string) error {

	// Create a temporary file
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name()) // Delete out temporary file

	// Copy source contents to temporary file
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = f.Write(contents)
	if err != nil {
		return err
	}
	defer f.Close()

	_, undoDeferGoList, undoFordeferList, err := common.FindIllegalStatements(path, contents)
	if err != nil {
		return err
	}

	// Run formatting algorithm
	contents, err = format.Source(contents)
	if err != nil {
		return err
	}

	// For undoing:
	var defersFound int
	var fordefersFound int

	fs := token.NewFileSet()

	_f := fs.AddFile(path, fs.Base(), len(contents))
	var s scanner.Scanner
	s.Init(_f, contents, nil, scanner.ScanComments)

	var state state

	for {
		pos, tok, l := s.Scan()
		if tok == token.EOF {
			break
		}

		if tok == token.DEFER {
			state = initial

			defersFound++

			if len(undoDeferGoList) != 0 {
				if undoDeferGoList[0].UndoType == common.DeferGo && undoDeferGoList[0].Pos == defersFound {
					addGoStmt(&contents, int(pos))
					undoDeferGoList = append(undoDeferGoList[:0], undoDeferGoList[1:]...)
				}
			}
		} else if tok == token.IDENT && l == "f0" {
			state = foFound
		} else if tok == token.PERIOD && state == foFound {
			state = dotFound
		} else if tok == token.IDENT && l == "d5" && state == dotFound {
			state = dfFound
		} else if tok == token.LPAREN && state == dfFound {
			state = lBrakFound
		} else if tok == token.RPAREN && state == lBrakFound {
			state = rBrakFound
		} else if tok == token.SEMICOLON && state == rBrakFound {
			state = initial

			fordefersFound++

			if len(undoFordeferList) != 0 {
				if undoFordeferList[0].Pos == fordefersFound {
					if undoFordeferList[0].UndoType == common.Fordefer {
						// fordefer
						addFordefer(&contents, false, int(pos))
						undoFordeferList = append(undoFordeferList[:0], undoFordeferList[1:]...)
					} else if undoFordeferList[0].UndoType == common.FordeferGo {
						// fordefer go
						addFordefer(&contents, true, int(pos))
						undoFordeferList = append(undoFordeferList[:0], undoFordeferList[1:]...)
					}
				}
			}
		} else {
			state = initial
		}
	}

	// Safe to overwrite original file
	err = ioutil.WriteFile(path, contents, 0644)
	if err != nil {
		return err
	}

	return nil
}

// addGoStmt inserts "go" statement after "defer".
func addGoStmt(data *[]byte, pos int) {
	str := " go"

	pos = pos - 1 + 5
	for i := 0; i < len(str); i++ {
		insert(data, pos+i, str[i])
	}
}

// addFordefer inserts "fordefer" and "fordefer go" statements
// after encountering a "f0.d5();" statement.
func addFordefer(data *[]byte, goStmt bool, pos int) {
	str := "fordefer "
	goStr := "go "

	pos = pos - 8

	// Replace f0.d5(); with fordefer
	for i := 0; i < len(str); i++ {
		(*data)[pos+i] = str[i]
	}

	pos = pos + len(str)
	if goStmt {
		for i := 0; i < len(goStr); i++ {
			insert(data, pos+i, goStr[i])
		}
		pos = pos + len(goStr)
	}

	// Remove indenting before next statement on next line
	for {
		if pos >= len((*data)) {
			break
		}

		if charFound := (*data)[pos]; charFound == []byte("	")[0] {
			(*data) = append((*data)[:pos], (*data)[pos+1:]...)
			continue
		} else {
			break
		}

		pos = pos + 1
	}
}

func insert(data *[]byte, idx int, char byte) {
	*data = append(*data, 0)
	copy((*data)[idx+1:], (*data)[idx:])
	(*data)[idx] = char
}
