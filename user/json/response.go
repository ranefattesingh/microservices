package json

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ranefattesingh/microservices/user/pkg"
)

var (
	ErrInternalServerErrorMsg = "server is unable to process at a moment"
)

type Responder interface {
	ResponseJSON(any) error
	Created(any) error
	NoContent() error
	BadRequest(error) error
	InternalServerError() error
	NotFound(error) error
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
	Fields  []pkg.FieldError `json:"fields"`
}

func Respond(w http.ResponseWriter) Responder {
	return &Response{
		w: w,
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

func (r *Response) NoContent() error {
	r.statusCode = http.StatusNoContent

	return r.encodeJSON()
}

func (r *Response) BadRequest(err error) error {
	r.statusCode = http.StatusBadRequest
	r.mapError(err)

	return r.encodeJSON()
}

func (r *Response) InternalServerError() error {
	r.statusCode = http.StatusInternalServerError
	r.Error = &Error{
		Message: ErrInternalServerErrorMsg,
	}

	return r.encodeJSON()
}

func (r *Response) NotFound(err error) error {
	r.statusCode = http.StatusNotFound
	r.mapError(err)

	return r.encodeJSON()
}

func (r *Response) mapError(err error) {
	r.Error = &Error{
		Message: err.Error(),
	}

	httpErr, ok := errors.AsType[*pkg.HTTPError](err)
	if ok {
		r.Error.Fields = httpErr.Fields
		r.Error.Message = httpErr.Message
	}

}

func (r *Response) encodeJSON() error {
	r.w.Header().Add("Content-Type", "application/json")
	r.w.WriteHeader(r.statusCode)
	err := json.NewEncoder(r.w).Encode(r)
	if err != nil {
		return err
	}

	return nil
}
