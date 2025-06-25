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

	ErrInvalidPath      = errors.New("invalid path")
	ErrFileNotFound     = errors.New("file not found")
	ErrFolderNotFound   = errors.New("folder not found")
	ErrFileExists       = errors.New("file already exists")
	ErrFolderExists     = errors.New("folder already exists")
	ErrInvalidFileType  = errors.New("invalid file type")
	ErrFolderNotEmpty   = errors.New("folder not empty")
	ErrFileTooLarge     = errors.New("file too large")
	ErrStorageError     = errors.New("storage error")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidOperation = errors.New("invalid operation")
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

func InvalidPath(msg string) *AppError {
	return NewAppError(http.StatusBadRequest, ErrInvalidPath, msg)
}

func FileNotFound(msg string) *AppError {
	return NewAppError(http.StatusNotFound, ErrFileNotFound, msg)
}

func FolderNotFound(msg string) *AppError {
	return NewAppError(http.StatusNotFound, ErrFolderNotFound, msg)
}

func FileExists(msg string) *AppError {
	return NewAppError(http.StatusConflict, ErrFileExists, msg)
}

func FolderExists(msg string) *AppError {
	return NewAppError(http.StatusConflict, ErrFolderExists, msg)
}

func InvalidFileType(msg string) *AppError {
	return NewAppError(http.StatusBadRequest, ErrInvalidFileType, msg)
}

func FolderNotEmpty(msg string) *AppError {
	return NewAppError(http.StatusConflict, ErrFolderNotEmpty, msg)
}

func FileTooLarge(msg string) *AppError {
	return NewAppError(http.StatusRequestEntityTooLarge, ErrFileTooLarge, msg)
}

func StorageError(msg string) *AppError {
	return NewAppError(http.StatusInternalServerError, ErrStorageError, msg)
}

func PermissionDenied(msg string) *AppError {
	return NewAppError(http.StatusForbidden, ErrPermissionDenied, msg)
}

func InvalidOperation(msg string) *AppError {
	return NewAppError(http.StatusBadRequest, ErrInvalidOperation, msg)
}

func TooManyRequests(msg string) *AppError {
	return NewAppError(http.StatusTooManyRequests, errors.New("too many requests"), msg)
}
