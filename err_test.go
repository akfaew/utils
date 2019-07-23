package utils

import (
	"errors"
	"strings"
	"testing"

	"github.com/akfaew/test"
)

func TestErrorc(t *testing.T) {
	test.True(t, Errorc(nil) == nil)

	e := Errorc(errors.New("oups")).Error()
	test.True(t, strings.HasSuffix(e, "utils/err_test.go:14 oups"))

	// two different errors with the same message
	e = Errorc(errors.New("oups"), errors.New("oups")).Error()
	test.True(t, strings.HasSuffix(e, "utils/err_test.go:18 oups"))

	// two identical errors
	ee := errors.New("oups")
	test.EqualStr(t, Errorc(ee, ee).Error(), "oups")
}
