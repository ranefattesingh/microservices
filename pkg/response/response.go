package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	httperror "github.com/ranefattesingh/ecommerce-platform/pkg/http_error"
)

// Response is the standard API envelope.
type Response struct {
	Success bool  `json:"success"`
	Data    any   `json:"data,omitempty"`
	Error   error `json:"error,omitempty"`
	Meta    *Meta `json:"meta,omitempty"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"totalPages,omitempty"`
}

func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func NoContent(c *gin.Context, data any) {
	c.JSON(http.StatusNoContent, nil)
}

// ErrorResponse sends anappropriate error with status code.
func ErrorResponse(c *gin.Context, err error) {
	httpErr := new(httperror.ErrorInfo)

	if errors.As(err, &httpErr) {
		switch httpErr.HTTPStatusCode {
		case http.StatusBadRequest:
			BadRequest(c, err)
		case http.StatusConflict:
			Conflict(c, err)
		default:
			InternalServerError(c)
		}
	}
}

// BadRequest sends an bad request error response.
func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error:   err,
	})
}

// Conflict sends an bad request error response.
func Conflict(c *gin.Context, err error) {
	c.JSON(http.StatusConflict, Response{
		Success: false,
		Error:   err,
	})
}

// InternalServerError sends an error response.
func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error:   errors.New("server fail to handle the request"),
	})
}
