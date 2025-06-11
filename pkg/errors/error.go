package errors

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const FormatDateTime = time.RFC3339

type ErrorResponse struct {
	Err     string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
	Time    string `json:"time,omitempty"`
}

type AppError struct {
	Code    int
	Err     error
	Message string
	Time    string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, error: %v, message: %s, time: %s", e.Code, e.Err, e.Message, e.Time)
}

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrEmptyParameter = errors.New("empty parameter")
	ErrDatabase       = errors.New("database error")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")
	ErrInternalServer = errors.New("internal server error")
)

func NewAppError(code int, err error, msg string) *AppError {
	return &AppError{
		Code:    code,
		Err:     err,
		Message: msg,
		Time:    time.Now().Format(FormatDateTime),
	}
}

func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, ErrorResponse{
			Err:     appErr.Err.Error(),
			Message: appErr.Message,
			Code:    appErr.Code,
			Time:    appErr.Time,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Err:     "internal_error",
		Message: "Unexpected error occurred",
		Code:    http.StatusInternalServerError,
		Time:    time.Now().Format(FormatDateTime),
	})
}

func BadRequest(msg string) *AppError {
	return NewAppError(http.StatusBadRequest, ErrInvalidRequest, msg)
}

func EmptyField(field string) *AppError {
	return NewAppError(http.StatusBadRequest, ErrEmptyParameter, fmt.Sprintf("Field '%s' is required", field))
}

func Database(msg string) *AppError {
	return NewAppError(http.StatusInternalServerError, ErrDatabase, msg)
}

func Unauthorized(msg string) *AppError {
	return NewAppError(http.StatusUnauthorized, ErrUnauthorized, msg)
}

func Forbidden(msg string) *AppError {
	return NewAppError(http.StatusForbidden, ErrForbidden, msg)
}

func NotFound(msg string) *AppError {
	return NewAppError(http.StatusNotFound, ErrNotFound, msg)
}

func Conflict(msg string) *AppError {
	return NewAppError(http.StatusConflict, ErrConflict, msg)
}

func InternalServer(msg string) *AppError {
	return NewAppError(http.StatusInternalServerError, ErrInternalServer, msg)
}
