package httperror

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ranefattesingh/ecommerce-platform/pkg/response"
)

var reservedKeys = map[string]struct{}{
	"type":     {},
	"title":    {},
	"status":   {},
	"detail":   {},
	"instance": {},
	// "violations":    {},
	// "requestId":     {},
	// "correlationId": {},
	// "traceId":       {},
	// "sampled":       {},
}

type ProblemDetails struct {
	Type       string
	Title      string
	Detail     string
	Instance   string
	Status     int
	Extensions map[string]any

	header http.Header
	err    error
}

var _ response.Responder = (*ProblemDetails)(nil)

func NewProblemDetails(title string, Status int) *ProblemDetails {
	return &ProblemDetails{
		Type:       "about:blank",
		Title:      title,
		Detail:     "",
		Instance:   "",
		Status:     Status,
		Extensions: make(map[string]any),
		err:        nil,
	}
}

func (p *ProblemDetails) SetType(t string) *ProblemDetails {
	p.Type = t
	return p
}

func (p *ProblemDetails) SetTitle(title string) *ProblemDetails {
	p.Title = title
	return p
}

func (p *ProblemDetails) SetDetail(detail string) *ProblemDetails {
	p.Detail = detail
	return p
}

func (p *ProblemDetails) SetInstance(instance string) *ProblemDetails {
	p.Instance = instance
	return p
}

func (p *ProblemDetails) SetStatus(code int) *ProblemDetails {
	p.Status = code
	return p
}

func (p *ProblemDetails) AddExtension(key string, value any) *ProblemDetails {
	if p.Extensions == nil {
		p.Extensions = make(map[string]any)
	}

	_, ok := reservedKeys[key]
	if ok {
		return p
	}

	p.Extensions[key] = value
	return p
}

func (p *ProblemDetails) SetError(err error) *ProblemDetails {
	p.err = errors.Join(p.err, err)
	return p
}

func (p *ProblemDetails) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')

	writeKV := func(key string, val any, first *bool) error {
		if !*first {
			buf.WriteByte(',')
		}

		kb, err := json.Marshal(key)
		if err != nil {
			return err
		}

		vb, err := json.Marshal(val)
		if err != nil {
			return err
		}

		buf.Write(kb)
		buf.WriteByte(':')
		buf.Write(vb)
		*first = false
		return nil
	}

	first := true

	if err := writeKV("type", p.Type, &first); err != nil {
		return nil, err
	}

	if err := writeKV("title", p.Title, &first); err != nil {
		return nil, err
	}

	if err := writeKV("detail", p.Detail, &first); err != nil {
		return nil, err
	}

	if err := writeKV("instance", p.Instance, &first); err != nil {
		return nil, err
	}

	if err := writeKV("status", p.Status, &first); err != nil {
		return nil, err
	}

	if len(p.Extensions) > 0 {
		keys := make([]string, 0, len(p.Extensions))
		for k := range p.Extensions {
			keys = append(keys, k)
		}

		slices.Sort(keys)
		for _, k := range keys {
			if err := writeKV(k, p.Extensions[k], &first); err != nil {
				return nil, err
			}
		}
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (p *ProblemDetails) Error() string {
	if p.err != nil {
		return p.err.Error()
	}

	return fmt.Sprintf("problem: %s: %s", p.Title, p.Detail)
}

func (p *ProblemDetails) Respond(c *gin.Context) error {
	if p.Status == 0 {
		p.Status = http.StatusInternalServerError
	}

	if p.Instance == "" {
		p.Instance = c.Request.URL.Path
	}

	if p.Type == "" {
		p.Type = "about:blank"
	}

	for k, values := range p.header {
		for _, value := range values {
			c.Writer.Header().Add(k, value)
		}
	}

	c.JSON(p.Status, p)

	return nil
}

func BadRequest(violations ...Violations) *ProblemDetails {
	v := mergeViolations(violations...)

	detail := "Request cannot be processed because the request was not understood by the server or malformed"
	if len(v) != 0 {
		detail = "Request validation failed. See 'violations' for details"
	}

	problemDetails := NewProblemDetails("Bad Request", http.StatusBadRequest).SetDetail(detail)

	if len(v) > 0 {
		problemDetails.AddExtension("violations", v)
	}

	return problemDetails
}

func Unauthorized() *ProblemDetails {
	return NewProblemDetails("Unauthorized", http.StatusUnauthorized).SetDetail("Authorization is required to perform this action.")
}

func Forbidden() *ProblemDetails {
	return NewProblemDetails("Forbidden", http.StatusForbidden).SetDetail("You don't have the required permission or role to perform this action.")
}

func NotFound() *ProblemDetails {
	return NewProblemDetails("Not Found", http.StatusNotFound).SetDetail("The requested resource could not be found.")
}

func MethodNotAllowed(allowedMethods ...string) *ProblemDetails {
	problemDetails := NewProblemDetails("Method Not Allowed", http.StatusMethodNotAllowed).SetDetail("The requested method is not allowed on this resource.")

	if len(allowedMethods) > 0 {
		problemDetails.header = http.Header{
			"Allow": []string{strings.Join(allowedMethods, ", ")},
		}
	}

	return problemDetails
}

func NotAcceptable() *ProblemDetails {
	return NewProblemDetails("Not Acceptable", http.StatusNotAcceptable).SetDetail("Server cannot produce a response matching the acceptable content types.")
}

func Conflict(violations ...Violations) *ProblemDetails {
	v := mergeViolations(violations...)
	detail := "The request could not be completed due to a conflict."
	if len(v) != 0 {
		detail = "The request could not be completed due to a conflict with existing resource. See 'violations' for details"
	}

	problemDetails := NewProblemDetails("Conflict", http.StatusConflict).SetDetail(detail)

	if len(v) > 0 {
		problemDetails.AddExtension("violations", v)
	}

	return problemDetails
}

func Gone() *ProblemDetails {
	return NewProblemDetails("Gone", http.StatusGone).SetDetail("The requested resource is no longer available.")
}

func UnsupportedMediaType() *ProblemDetails {
	return NewProblemDetails("Unsupported Media Type", http.StatusUnsupportedMediaType).SetDetail("The request entity has media type which the server or resource does not support.")
}

func UnprocessableEntity(violations ...Violations) *ProblemDetails {
	v := mergeViolations(violations...)
	detail := "The request entity could not be processed."
	if len(v) != 0 {
		detail = "The request entity could not be processed due to validation errors. See 'violations' for details"
	}
	problemDetails := NewProblemDetails("Unprocessable Entity", http.StatusUnprocessableEntity).SetDetail(detail)

	if len(v) > 0 {
		problemDetails.AddExtension("violations", v)
	}

	return problemDetails
}

func TooManyRequests() *ProblemDetails {
	return NewProblemDetails("Too Many Requests", http.StatusTooManyRequests).SetDetail("The request was rate-limited. Try again later.")
}

func InternalServerError(errs ...error) *ProblemDetails {
	problemDetails := NewProblemDetails("Internal Server Error", http.StatusInternalServerError).SetDetail("Server encountered an internal error processing the request. Please try again later")

	if len(errs) > 0 {
		err := errors.Join(errs...)
		problemDetails.err = err
	}

	return problemDetails
}

func BadGateway() *ProblemDetails {
	return NewProblemDetails("Bad Gateway", http.StatusBadGateway).SetDetail("The server received an invalid response from an upstream server..")
}

func ServiceUnavailable() *ProblemDetails {
	return NewProblemDetails("Service Unavailable", http.StatusServiceUnavailable).SetDetail("The server is currently unable to handle the request due to a temporary overloading or maintenance of the server. Please try again later")
}

type Violations map[string][]string

func (v *Violations) Add(field string, errors ...string) {
	if *v == nil {
		*v = make(Violations)
	}

	m := *v
	violation, ok := m[field]
	if !ok {
		violation = make([]string, 0, len(errors))
	}

	m[field] = append(violation, errors...)
}

func (v *Violations) Len() int {
	if v == nil {
		return 0
	}

	return len(*v)
}

func (v Violations) AsResponse() *ProblemDetails {
	if len(v) == 0 {
		return UnprocessableEntity()
	}

	return UnprocessableEntity(v)
}

func mergeViolations(v ...Violations) Violations {
	merged := make(Violations)
	for _, violation := range v {
		for field, msgs := range violation {
			merged[field] = append(merged[field], msgs...)
		}
	}

	return merged
}
