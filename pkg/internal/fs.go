package internal

import (
	"path/filepath"
	"runtime"
)

// GetPackagePath returns the directory of the package currently running this program.
// Bare in mind it's not the same as the CWD of the bin. Returns empty string if
// the program running does not following the pkg dir format
func GetPackagePath() string {
	_, sourceCode, _, _ := runtime.Caller(0)
	for dir, last := filepath.Split(sourceCode); dir != ""; dir, last = filepath.Split(filepath.Clean(dir)) {
		if last == "pkg" {
			return dir
		}
	}

	return ""
}
