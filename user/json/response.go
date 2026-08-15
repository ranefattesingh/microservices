package json

import (
	"encoding/json"
	"net/http"

	"github.com/ranefattesingh/microservices/user/pkg"
)

const ErrInternalServerErrorMsg = "server is unable to process the request at the moment"

type Responder interface {
	JSON(any) error
	Created(any) error
	NoContent() error
	BadRequest(error, ...pkg.FieldError) error
	InternalServerError() error
	NotFound(error) error
	Conflict(err error) error
}

type Response struct {
	Success    bool   `json:"success"`
	Data       any    `json:"data,omitempty"`
	Error      *Error `json:"error,omitempty"`
	statusCode int    `json:"-"`
	w          http.ResponseWriter
}

type Error struct {
	Message string           `json:"message"`
	Fields  []pkg.FieldError `json:"fields,omitempty"`
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

func (r *Response) NoContent() error {
	r.w.WriteHeader(http.StatusNoContent)
	return nil
}

func (r *Response) BadRequest(err error, fieldErrs ...pkg.FieldError) error {
	r.Success = false
	r.Data = nil
	r.statusCode = http.StatusBadRequest
	r.Error = &Error{
		Message: err.Error(),
		Fields:  fieldErrs,
	}

	return r.encodeJSON()
}

func (r *Response) NotFound(err error) error {
	r.Success = false
	r.Data = nil
	r.statusCode = http.StatusNotFound
	r.Error = &Error{
		Message: err.Error(),
	}

	return r.encodeJSON()
}

func (r *Response) Conflict(err error) error {
	r.Success = false
	r.Data = nil
	r.statusCode = http.StatusConflict
	r.Error = &Error{
		Message: err.Error(),
	}

	return r.encodeJSON()
}

func (r *Response) InternalServerError() error {
	r.Success = false
	r.Data = nil
	r.statusCode = http.StatusInternalServerError
	r.Error = &Error{
		Message: ErrInternalServerErrorMsg,
	}

	return r.encodeJSON()
}

func (r *Response) encodeJSON() error {
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(r.statusCode)
	return json.NewEncoder(r.w).Encode(r)
}
