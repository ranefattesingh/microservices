package json

import (
	"encoding/json"
	"net/http"
)

type Responder interface {
	JSON(any) error
	Created(any) error
	BadRequest(error, ...FieldError) error
	InternalServerError() error
}

type Response struct {
	Success    bool        `json:"success"`
	Data       any         `json:"data,omitempty"`
	Error      *Error      `json:"error,omitempty"`
	statusCode int         `json:"-"`
	w          http.ResponseWriter
}

type Error struct {
	Message string       `json:"message"`
	Fields  []FieldError `json:"fields,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"error"`
}

func Respond(w http.ResponseWriter) Responder {
	return &Response{
		w:          w,
		statusCode: http.StatusOK,
	}
}

func (r *Response) JSON(data any) error {
	r.Data = data
	r.Success = true
	r.Error = nil
	r.statusCode = http.StatusOK

	return r.encodeJSON()
}

func (r *Response) Created(data any) error {
	r.Data = data
	r.Success = true
	r.Error = nil
	r.statusCode = http.StatusCreated

	return r.encodeJSON()
}

func (r *Response) BadRequest(err error, fieldErrs ...FieldError) error {
	r.Success = false
	r.Data = nil
	r.statusCode = http.StatusBadRequest
	r.Error = &Error{
		Message: err.Error(),
		Fields:  fieldErrs,
	}

	return r.encodeJSON()
}

func (r *Response) InternalServerError() error {
	r.Success = false
	r.Data = nil
	r.statusCode = http.StatusInternalServerError
	r.Error = &Error{
		Message: "server is unable to process the request at the moment",
	}

	return r.encodeJSON()
}

func (r *Response) encodeJSON() error {
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(r.statusCode)
	return json.NewEncoder(r.w).Encode(r)
}
