package service

type Code int

const (
	Validation Code = iota
	NotFound
	Conflict
	Unauthorized
	Forbidden
	Internal
)

type FieldError struct {
	Field   string
	Message string
}

type Error struct {
	Code       Code
	Message    string
	Violations []FieldError
}

func (e *Error) Error() string {
	return e.Message
}
