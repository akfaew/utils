package utils

import (
	"path/filepath"
	"runtime"
	"strings"
)

var trimprefix = ""

// Error messages get the file path prepended. Let's skip the full path as it's user facing, and just keep the filename.
func Init(path string) {
	if path == "" {
		_, path, _, _ = runtime.Caller(1)
	}

	trimprefix = filepath.Dir(path) + string(filepath.Separator)
}

func trim(path string) string {
	return strings.TrimPrefix(path, trimprefix)
}

func logctx(skip int) (file string, line int) {
	_, file, line, _ = runtime.Caller(skip + 1)
	file = trim(file)

	return
}
