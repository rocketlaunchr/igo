// Copyright 2018 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/rocketlaunchr/igo/cmds"
)

// These are the files we want to transpile
var sourceFiles []string

// tempGeneratedFiles maps the actual source file to the temporary generated file
var tempGeneratedFiles map[string]string

func init() {
	tempGeneratedFiles = make(map[string]string)

	sourceFiles = []string{
		"test/main.igo",
	}
}

func main() {

	var rootCmd = &cobra.Command{
		Use: "igo",
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version prints the igo version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("igo version: 0.0.1")
		},
	}

	var buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Transpile igo files to go files",
		Run:   cmds.BuildCmd,
	}

	rootCmd.AddCommand(buildCmd, versionCmd)
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
