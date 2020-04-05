// Copyright 2018-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package file

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"go/format"
	"go/token"
	"io/ioutil"
)

// ReadFile will read a source file and return the contents as a
// []byte.
func ReadFile(sourceFile string) ([]byte, error) {
	return ioutil.ReadFile(sourceFile)
}

// CreateTempFile creates a copy of the source file and
// returns the new destination file name.
func CreateTempFile(basePath string, sourceBytes []byte) (string, error) {

	// Determine file's md5 hash
	destinationFile := basePath + fmt.Sprintf("i%s.go", md5OfData(sourceBytes))

	// Copy source file to destination file
	err := ioutil.WriteFile(destinationFile, sourceBytes, 0644)
	if err != nil {
		return "", err
	}

	return destinationFile, nil
}

func md5OfData(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

// UpdateFile will replace an existing file with new data.
func UpdateFile(fileName string, data []byte) error {

	err := ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func SaveFmtFile(fileName string, fset *token.FileSet, node interface{}) error {

	var buf bytes.Buffer
	err := format.Node(&buf, fset, node)
	if err != nil {
		return err
	}

	return UpdateFile(fileName, buf.Bytes())
}

func PrependFile(fileName string, header string) error {
	b, err := ReadFile(fileName)
	if err != nil {
		return err
	}

	err = UpdateFile(fileName, append([]byte(header), b...))
	if err != nil {
		return err
	}

	return nil
}
