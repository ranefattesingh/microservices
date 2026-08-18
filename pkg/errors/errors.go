package errors

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"error"`
}

type HTTPError struct {
	Message string       `json:"error"`
	Fields  []FieldError `json:"fields,omitempty"`
}

func NewHTTPError(msg string, fields ...FieldError) *HTTPError {
	return &HTTPError{
		Message: msg,
		Fields:  fields,
	}
}

func (e *HTTPError) Error() string {
	return e.Message
}
