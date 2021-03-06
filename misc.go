package utils

import (
	"fmt"
	"hash/crc32"
	"reflect"
	"strings"
	"sync"
)

func Slash(text string) (string, string) {
	res := strings.SplitN(text, "/", 2)
	if len(res) == 1 {
		return res[0], ""
	}
	return res[0], res[1]
}

func HasElem(arr interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(arr)

	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			// panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}

	return false
}

type ErrorList struct {
	sync.Mutex
	errs []error
}

func (errs *ErrorList) Append(err error) {
	errs.Lock()
	defer errs.Unlock()
	errs.errs = append(errs.errs, err)
}

func (errs *ErrorList) Error() error {
	errs.Lock()
	defer errs.Unlock()

	return Errors(errs.errs)
}

func Errors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var errstrings []string
	for _, e := range errs {
		if e != nil {
			errstrings = append(errstrings, e.Error())
		}
	}

	return fmt.Errorf("[\"%s\"]", strings.Join(errstrings, "\", \""))
}

// A simple sum for naming fixture files in tests, e.g. based on an URL.
func Sum(txt string) string {
	return fmt.Sprintf("%08x", crc32.Checksum([]byte(txt), crc32.IEEETable))
}
