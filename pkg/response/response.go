package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Responder interface {
	Respond(c *gin.Context) error
}

type ResponseBuilder struct {
	status  int
	header  http.Header
	cookies []*http.Cookie
	body    any
}

var _ Responder = (*ResponseBuilder)(nil)

func Response() *ResponseBuilder {
	return &ResponseBuilder{
		status:  http.StatusOK,
		header:  http.Header{},
		cookies: make([]*http.Cookie, 0),
	}
}

func (rb *ResponseBuilder) Status(status int) *ResponseBuilder {
	rb.status = status
	return rb
}

func (rb *ResponseBuilder) Header(key string, values ...string) *ResponseBuilder {
	for _, value := range values {
		rb.header.Add(key, value)
	}

	return rb
}

func (rb *ResponseBuilder) Headers(headers http.Header) *ResponseBuilder {
	for key, values := range headers {
		for _, value := range values {
			rb.header.Add(key, value)
		}
	}

	return rb
}

func (rb *ResponseBuilder) Cookie(cookie *http.Cookie) *ResponseBuilder {
	rb.cookies = append(rb.cookies, cookie)
	return rb
}

func (rb *ResponseBuilder) Cookies(cookie []*http.Cookie) *ResponseBuilder {
	rb.cookies = append(rb.cookies, cookie...)
	return rb
}

func (rb *ResponseBuilder) Body(body any) *ResponseBuilder {
	rb.body = body
	return rb
}

func (rb *ResponseBuilder) ContentType(contentType string) *ResponseBuilder {
	rb.header.Set("Content-Type", contentType)
	return rb
}

func (rb *ResponseBuilder) Respond(c *gin.Context) error {
	for key, values := range rb.header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	for _, cookie := range rb.cookies {
		http.SetCookie(c.Writer, cookie)
	}

	if rb.body == nil {
		c.Status(rb.status)
		return nil
	}

	switch body := rb.body.(type) {
	case []byte:
		c.Data(rb.status, "application/octet-stream", body)
	case string:
		c.String(rb.status, body)
	default:
		c.JSON(rb.status, body)
	}

	return nil
}

func Ok(body any) *ResponseBuilder {
	return Response().Status(http.StatusOK).Body(body)
}

func Created(location string) *ResponseBuilder {
	return Response().Status(http.StatusCreated).Header("Location", location)
}

func Accepted() *ResponseBuilder {
	return Response().Status(http.StatusAccepted)
}

func NoContent() *ResponseBuilder {
	return Response().Status(http.StatusNoContent)
}

func HTML(body string) *ResponseBuilder {
	return Response().
		Header("Content-Type", "text/html; charset=utf-8").
		Body(body)
}

func JSON(status int, body any) *ResponseBuilder {
	return Response().
		Status(status).
		Body(body)
}
