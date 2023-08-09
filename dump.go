package utils

import (
	"encoding/json"
	"os"
)

func Dump(path string, what interface{}) error {
	if b, err := json.MarshalIndent(what, "", "\t"); err != nil {
		return Errorc(err)
	} else if err := os.WriteFile(path, b, 0644); err != nil {
		return Errorc(err)
	}

	return nil
}
