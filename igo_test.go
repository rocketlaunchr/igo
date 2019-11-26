package main

import (
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

	expByte := make([]byte, chunkSize)
	_, err = exp.Read(expByte)
	if err != nil {
		return false, err
	}

	genByte := make([]byte, chunkSize)
	_, err = gen.Read(genByte)
	if err != nil {
		return false, err
	}

	// Comparing files using regexp patterns from expected_code file
	matched, err := regexp.MatchString(string(expByte), string(genByte))
	if err != nil {
		return false, err
	}

	return matched, nil
}

func TestIgoBuild(t *testing.T) {
	igoFile := "./test_files/test_sample.igo"
	genFile := "./test_files/gen_test_sample.go"
	expCode := "./test_files/expected.go"

	// Generating go code from igo sample file
	cmds.BuildCmd(buildCmd, []string{igoFile})

	confirm, err := fileCompare(expCode, genFile)
	if err != nil {
		t.Errorf("Test failed with error: %s\n", err)
	}

	if !confirm {
		t.Error("Generated Code is not same with expected code.")
	}

	// Delete generated file
	if err := os.Remove(genFile); err != nil {
		t.Errorf("Test failed with error:%s\n", err)
	}

}
