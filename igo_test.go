package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rocketlaunchr/igo/cmds"
	"github.com/spf13/cobra"
	"io"
	"os"
	"testing"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Transpile igo files to go files",
	Run:   cmds.BuildCmd,
}

const chunkSize = 64000

func fileCompare(file1, file2 string) (bool, error) {
	// Check file size ...

	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF { // end of both files reached
				return true, nil
			} else if err1 == io.EOF || err2 == io.EOF {
				return false, nil
			} else {
				return false, errors.New(fmt.Sprintf("Errors encountered from the two files.\nFile1:\n%s\nFile2:\n%s\n", err1, err2))
			}
		}

		if !bytes.Equal(b1, b2) {
			return false, nil
		}
	}
}

func TestIgo(t *testing.T) {
	igoFile := "./test_files/test_sample.igo"
	genFile := "./test_files/gen_test_sample.go"
	expCode := "./test_files/expected_build_code"

	// Generating go code from igo sample file
	cmds.BuildCmd(buildCmd, []string{igoFile})

	t.Logf("Comparing generated file code '%s' and expected build code file '%s'\n", genFile, expCode)
	confirm, err := fileCompare(expCode, genFile)
	if err != nil {
		t.Errorf("Test failed with error: %s\n", err)
	}
	t.Log("Done Comparing files.\n")
	if !confirm {
		t.Errorf("Generated Code is not same with expected code.")
	}

	// delete generated file
	t.Log("deleting Generated Test_build File...\n")
	if err := os.Remove(genFile); err != nil {
		t.Errorf("Test failed with error:%s\n", err)
	}
	t.Log("deleted.\n")

}
