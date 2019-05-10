package utils

import (
	"testing"

	"github.com/akfaew/test"
)

func Test_Init(t *testing.T) {
	Init("/a/b/c/d/e.go")
	test.EqualStr(t, "/a/b/c/d/", trimprefix)

	Init("")
	file, _ := logctx(0)
	test.EqualStr(t, file, "utils_test.go")
}
