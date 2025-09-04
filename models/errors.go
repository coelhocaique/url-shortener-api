package models

import (
	"errors"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Message    string
	StatusCode int
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// GetStatusCode returns the HTTP status code for the error
func (e *AppError) GetStatusCode() int {
	return e.StatusCode
}

// Custom error instances
var (
	ErrInvalidURLFormat     = &AppError{Message: "invalid URL format", StatusCode: http.StatusBadRequest}
	ErrInvalidURLScheme     = &AppError{Message: "invalid URL: missing scheme or host", StatusCode: http.StatusBadRequest}
	ErrInvalidAliasLength   = &AppError{Message: "alias must be between 3 and 20 characters", StatusCode: http.StatusBadRequest}
	ErrInvalidAliasChars    = &AppError{Message: "alias can only contain letters, numbers, and hyphens", StatusCode: http.StatusBadRequest}
	ErrAliasAlreadyExists   = &AppError{Message: "alias already exists", StatusCode: http.StatusConflict}
	ErrShortCodeNotFound    = &AppError{Message: "short code not found", StatusCode: http.StatusNotFound}
	ErrShortCodeExpired     = &AppError{Message: "short code has expired", StatusCode: http.StatusNotFound}
)

// GetStatusCodeFromError extracts HTTP status code from an error
func GetStatusCodeFromError(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.GetStatusCode()
	}
	// Default to 500 for unknown errors
	return http.StatusInternalServerError
}
