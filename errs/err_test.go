package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	myErr := fmt.Errorf("my error")

	t.Run("nil", func(t *testing.T) {
		assert.NoError(t, Wrap(nil))
		assert.NoError(t, Wrapf(nil, "something failed"))
	})

	t.Run("errors.Is()", func(t *testing.T) {
		assert.True(t, errors.Is(Wrap(myErr), myErr))
		assert.True(t, errors.Is(Wrapf(myErr, "something failed"), myErr))
	})

	t.Run("paths", func(t *testing.T) {
		thisfile := "/utils/errs/err_test.go:"
		assert.Contains(t, Wrap(myErr).Error(), thisfile)
		assert.Contains(t, Wrapf(myErr, "something failed").Error(), thisfile)
	})
}
