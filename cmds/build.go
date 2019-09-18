// Copyright 2018-19 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package cmds

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rocketlaunchr/igo/addressable"
	"github.com/rocketlaunchr/igo/builtin"
	"github.com/rocketlaunchr/igo/common"
	"github.com/rocketlaunchr/igo/config"
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
			path, filename := filepath.Split(igoFile)
			newFileName := path + "gen_" + strings.TrimSuffix(filename, "igo") + "go"
			err := os.Rename(goFile, newFileName)
			if err != nil {
				os.Exit(common.Report(err))
				return
			}

			// Add igo header
			header := `// DO NOT MODIFY! AUTO GENERATED BY igo v%s (https://github.com/rocketlaunchr/igo)
`

			header = fmt.Sprintf(header, config.VERSION)
			err = file.PrependFile(newFileName, header)
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

	// Cleanup nuisance comments (inside functions which interfere with modifying tree)
	err = common.RemoveComments(tempFileName)
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

	err = builtin.Process(tempFileName, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
