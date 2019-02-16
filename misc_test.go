package utils

import (
	"fmt"
	"syscall"
	"testing"

	"github.com/akfaew/test"
)

func TestSlash(t *testing.T) {
	tests := []struct {
		input string
		ret1  string
		ret2  string
	}{
		{"", "", ""},
		{"test", "test", ""},
		{"test/", "test", ""},
		{"test///", "test", "//"},
		{"test/case", "test", "case"},
		{"test/case/a/b/c", "test", "case/a/b/c"},
		{"/case/a/b/c", "", "case/a/b/c"},
	}

	for _, tc := range tests {
		ret1, ret2 := Slash(tc.input)
		test.True(t, ret1 == tc.ret1)
		test.True(t, ret2 == tc.ret2)
	}
}

func TestHasElem(t *testing.T) {
	tests := []struct {
		arr  interface{}
		elem interface{}
		ret  bool
	}{
		{"", "", false},
		{5, "", false},
		{[]int{}, 2, false},
		{[]int{1}, 1, true},
		{[]int{1}, 2, false},
		{[]int{1, 2}, 2, true},
		{[]int{1, 2, 3}, 2, true},
		{[]string{}, 2, false},
		{[]string{}, "", false},
		{[]string{"a"}, "", false},
		{[]string{"a"}, "a", true},
		{[]string{"a", "b"}, "b", true},
		{[]string{"a", "b", "c"}, "b", true},
	}

	for _, tc := range tests {
		ret := HasElem(tc.arr, tc.elem)
		test.True(t, ret == tc.ret)
	}
}

func TestErrors(t *testing.T) {
	errs := []error{
		fmt.Errorf("First error: server error"),
		fmt.Errorf("Second error: %s", syscall.ENOPKG.Error()),
		fmt.Errorf("Third error: %s", syscall.ENOTCONN.Error()),
	}
	test.Fixture(t, Errors(errs).Error())
}
