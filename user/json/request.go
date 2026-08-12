package json

import (
	"encoding/json"
	"net/http"
)

func Decode(r *http.Request, req any) error {
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	return nil
}
