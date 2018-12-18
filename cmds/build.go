// Copyright 2018 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package cmds

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rocketlaunchr/igo/addressable"
	"github.com/rocketlaunchr/igo/common"
	"github.com/rocketlaunchr/igo/defergo"
	"github.com/rocketlaunchr/igo/file"
	"github.com/rocketlaunchr/igo/fordefer"
)

// tempGeneratedFiles maps the actual source file to the temporary generated file
var tempGeneratedFiles map[string]string

func init() {
	tempGeneratedFiles = make(map[string]string)
}

func BuildCmd(cmd *cobra.Command, args []string) {

	for _, path := range args {

		files, err := common.Files(path, "igo")
		if err != nil {
			os.Exit(common.Report(err))
			return
		}

		for _, path := range files {
			err := processFile(path)
			if err != nil {
				os.Exit(common.Report(err))
				return
			}
		}

		// Rename temporary files to *.go files
		for igoFile, goFile := range tempGeneratedFiles {
			newFileName := strings.TrimSuffix(igoFile, "igo") + "go"
			err := os.Rename(goFile, newFileName)
			if err != nil {
				os.Exit(common.Report(err))
				return
			}
		}
	}
}

func processFile(sourceFile string) error {

	b, err := file.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	goPos, _, _, err := common.FindIllegalStatements(sourceFile, b)
	if err != nil {
		return err
	}

	path, _ := filepath.Split(sourceFile)

	// Create a temporary file after preprocessing
	tempFileName, err := file.CreateTempFile(path, b)
	if err != nil {
		return err
	}
	tempGeneratedFiles[sourceFile] = tempFileName

	// Update temp generated file
	err = defergo.Process(tempFileName, goPos)
	if err != nil {
		return err
	}

	err = fordefer.Process(tempFileName)
	if err != nil {
		return err
	}

	err = addressable.Process(tempFileName)
	if err != nil {
		return err
	}

	return nil
}
