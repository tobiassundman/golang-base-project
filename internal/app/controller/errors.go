package controller

import (
	"fmt"
	"net/http"

	"github.com/tobiassundman/go-demo-app/internal/app/service"
)

type APIError struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"error_message"`
	Status    int    `json:"status"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("error code: %s, message: %s, status: %d", e.ErrorCode, e.Message, e.Status)
}

var (
	ErrUserNotFound = &APIError{
		ErrorCode: "ErrUserNotFound",
		Message:   "user not found",
		Status:    http.StatusNotFound,
	}
	ErrUserAlreadyExists = &APIError{
		ErrorCode: "ErrUserAlreadyExists",
		Message:   "user already exists",
		Status:    http.StatusConflict,
	}
	ErrValidationFailed = &APIError{
		ErrorCode: "ErrValidationFailed",
		Message:   "validation failed",
		Status:    http.StatusBadRequest,
	}
	ErrInternalServer = &APIError{
		ErrorCode: "ErrInternalServer",
		Message:   "internal server error",
		Status:    http.StatusInternalServerError,
	}
	ErrInvalidID = &APIError{
		ErrorCode: "ErrInvalidID",
		Message:   "invalid id",
		Status:    http.StatusBadRequest,
	}
)

// apiErrorFromServiceError converts service errors to API errors.
func apiErrorFromServiceError(err error) *APIError {
	switch err {
	case service.ErrUserNotFound:
		return ErrUserNotFound
	case service.ErrUserAlreadyExists:
		return ErrUserAlreadyExists
	default:
		return ErrInternalServer
	}
}
