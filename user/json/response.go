package json

import (
	"encoding/json"
	"net/http"
)

type Responder interface {
	ResponseJSON(data any) error
	BadRequest(error) error
}

type Response struct {
	Success    bool  `json:"success"`
	Data       any   `json:"data,omitempty"`
	Error      error `json:"error,omitempty"`
	statusCode int   `json:"-"`
	w          http.ResponseWriter
}

func Respond(w http.ResponseWriter) Responder {
	return &Response{
		statusCode: http.StatusOK,
		w:          w,
	}
}

func (r *Response) ResponseJSON(data any) error {
	r.Data = data
	r.Success = true

	err := json.NewEncoder(r.w).Encode(r.Data)
	if err != nil {
		return err
	}

	return nil
}

func (r *Response) BadRequest(err error) error {
	r.Error = err

	newErr := json.NewEncoder(r.w).Encode(r.Data)
	if newErr != nil {
		return newErr
	}

	return nil
}
