package json

import (
	"encoding/json"
	"net/http"
)

func Decode(r *http.Request, req any) error {
	return json.NewDecoder(r.Body).Decode(req)
}
