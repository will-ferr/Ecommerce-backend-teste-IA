package services

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrValidationFailed   ErrorCode = "VALIDATION_FAILED"
	ErrNotFound           ErrorCode = "NOT_FOUND"
	ErrUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrForbidden          ErrorCode = "FORBIDDEN"
	ErrConflict           ErrorCode = "CONFLICT"
	ErrRateLimit          ErrorCode = "RATE_LIMIT"
	ErrInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrTimeout            ErrorCode = "TIMEOUT"
	ErrDatabaseError      ErrorCode = "DATABASE_ERROR"
	ErrCacheError         ErrorCode = "CACHE_ERROR"
)

type AppError struct {
	Code      ErrorCode   `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Cause     error       `json:"-"`
	Timestamp string      `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
	UserID    uint        `json:"user_id,omitempty"`
}

type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

func NewError(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Timestamp: getCurrentTimestamp(),
	}
}

func NewErrorWithDetails(code ErrorCode, message string, details interface{}) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: getCurrentTimestamp(),
	}
}

func NewValidationError(field, message string, value interface{}) *AppError {
	return &AppError{
		Code:    ErrValidationFailed,
		Message: "Validation failed",
		Details: []ValidationError{{
			Field:   field,
			Message: message,
			Value:   value,
		}},
		Timestamp: getCurrentTimestamp(),
	}
}

func NewErrorWithCause(code ErrorCode, message string, cause error) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Cause:     cause,
		Timestamp: getCurrentTimestamp(),
	}
}

func (e *AppError) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("%s: %s (details: %+v)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) WithRequestID(requestID string) *AppError {
	e.RequestID = requestID
	return e
}

func (e *AppError) WithUserID(userID uint) *AppError {
	e.UserID = userID
	return e
}

func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

func (e *AppError) ToHTTPStatus() int {
	switch e.Code {
	case ErrValidationFailed:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrConflict:
		return http.StatusConflict
	case ErrRateLimit:
		return http.StatusTooManyRequests
	case ErrInternalError, ErrDatabaseError, ErrCacheError:
		return http.StatusInternalServerError
	case ErrServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrTimeout:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}

func getCurrentTimestamp() string {
	return fmt.Sprintf("%d", getCurrentTimestampMillis())
}

func getCurrentTimestampMillis() int64 {
	return 0 // Placeholder - should use time.Now().UnixMilli()
}
