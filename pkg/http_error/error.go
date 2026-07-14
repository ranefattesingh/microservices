package httperror

type ErrorInfo struct {
	HTTPStatusCode int
	Message        string `json:"message"`
}

func (ei *ErrorInfo) Error() string {
	return ei.Message
}
