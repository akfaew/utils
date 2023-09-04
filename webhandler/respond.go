package webhandler

import (
	"encoding/json"
	"net/http"

	"github.com/akfaew/utils"
)

func WriteResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if data == nil {
		// Write has other side effects, and we don't want to write "null".
		_, err := w.Write([]byte{})
		return utils.Errorc(err)
	}

	ret, err := json.Marshal(data)
	if err != nil {
		return utils.Errorc(err)
	}
	if _, err := w.Write(ret); err != nil {
		return utils.Errorc(err)
	}

	return nil
}
