package json

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Responder interface {
	ResponseJSON(data any) error
	Created(data any) error
	BadRequest(error) error
	InternalServerError() error
	NotFound() error
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

func (r *Response) Created(data any) error {
	r.Data = data
	r.Success = true
	r.statusCode = http.StatusCreated

	return r.ResponseJSON(r.Data)
}

func (r *Response) BadRequest(err error) error {
	r.statusCode = http.StatusBadRequest
	r.Error = err

	newErr := json.NewEncoder(r.w).Encode(r.Data)
	if newErr != nil {
		return newErr
	}

	return nil
}

func (r *Response) InternalServerError() error {
	r.statusCode = http.StatusInternalServerError
	r.Error = errors.New("internal server error")

	newErr := json.NewEncoder(r.w).Encode(r.Data)
	if newErr != nil {
		return newErr
	}

	return nil
}

func (r *Response) NotFound() error {
	r.statusCode = http.StatusNotFound
	r.Error = errors.New("not found")

	newErr := json.NewEncoder(r.w).Encode(r.Data)
	if newErr != nil {
		return newErr
	}

	return nil
}
