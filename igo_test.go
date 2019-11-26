package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/rocketlaunchr/igo/cmds"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Transpile igo files to go files",
	Run:   cmds.BuildCmd,
}

const chunkSize = 64000

func fileCompare(expected, generated string) (bool, error) {

	exp, err := os.Open(expected)
	if err != nil {
		return false, err
	}
	defer exp.Close()

	gen, err := os.Open(generated)
	if err != nil {
		return false, err
	}
	defer gen.Close()

	for {
		expByte := make([]byte, chunkSize)
		_, err1 := exp.Read(expByte)

		genByte := make([]byte, chunkSize)
		_, err2 := gen.Read(genByte)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF { // end of both files reached
				return true, nil
			} else if err1 == io.EOF || err2 == io.EOF {
				return false, nil
			} else {
				return false, fmt.Errorf("errors encountered from the two files, file1: %s file2: %s", err1, err2)
			}
		}

		// Comparing files using regexp patterns from expected_code file

		matched, err := regexp.MatchString(string(expByte), string(genByte))
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}

	}
}

func TestIgoBuild(t *testing.T) {
	igoFile := "./test_files/test_sample.igo"
	genFile := "./test_files/gen_test_sample.go"
	expCode := "./test_files/expected.go"

	// Generating go code from igo sample file
	cmds.BuildCmd(buildCmd, []string{igoFile})

	t.Logf("Comparing generated file code '%s' and expected build code file '%s'\n", genFile, expCode)
	confirm, err := fileCompare(expCode, genFile)
	if err != nil {
		t.Errorf("Test failed with error: %s\n", err)
	}
	t.Log("Done Comparing files.\n")
	if !confirm {
		t.Error("Generated Code is not same with expected code.")
	}

	// delete generated file
	t.Log("deleting Generated Test_build File...\n")
	if err := os.Remove(genFile); err != nil {
		t.Errorf("Test failed with error:%s\n", err)
	}
	t.Log("generated file deleted.\n")

}
