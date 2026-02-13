package utils

import (
	"fmt"
	"net/http"
)

// AppError represents a structured application error.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewAppError creates a new AppError with the given code, message, and detail.
func NewAppError(code int, message, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// ErrNotFound returns a 404 Not Found error.
func ErrNotFound(detail string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: "리소스를 찾을 수 없습니다",
		Detail:  detail,
	}
}

// ErrUnauthorized returns a 401 Unauthorized error.
func ErrUnauthorized(detail string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "인증이 필요합니다",
		Detail:  detail,
	}
}

// ErrBadRequest returns a 400 Bad Request error.
func ErrBadRequest(detail string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "잘못된 요청입니다",
		Detail:  detail,
	}
}

// ErrForbidden returns a 403 Forbidden error.
func ErrForbidden(detail string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: "접근 권한이 없습니다",
		Detail:  detail,
	}
}

// ErrValidation returns a 422 Unprocessable Entity error for validation failures.
func ErrValidation(detail string) *AppError {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: "입력값 검증에 실패했습니다",
		Detail:  detail,
	}
}

// ErrInternal returns a 500 Internal Server Error.
func ErrInternal(detail string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "서버 내부 오류가 발생했습니다",
		Detail:  detail,
	}
}

// ErrConflict returns a 409 Conflict error.
func ErrConflict(detail string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: "리소스 충돌이 발생했습니다",
		Detail:  detail,
	}
}
