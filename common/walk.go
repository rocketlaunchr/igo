// Copyright 2018 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package common

import (
	"os"
	"path/filepath"
	"strings"
)

// Files accepts a path and returns a slice containing the path
// or paths (if path is a directory). acceptExt is used to filter
// specific file extensions.
func Files(path string, acceptExt string) ([]string, error) {

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		// Check file extension
		if checkFileExt(info, acceptExt) {
			return []string{path}, nil
		} else {
			return []string{}, nil
		}
	}

	// The path is a directory
	allFiles := []string{}
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && checkFileExt(info, acceptExt) {
			allFiles = append(allFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return allFiles, nil
}

func checkFileExt(f os.FileInfo, acceptExt string) bool {
	name := f.Name()
	ext := "." + acceptExt
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ext)
}
