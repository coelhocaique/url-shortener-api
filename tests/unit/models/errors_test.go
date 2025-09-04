package main

import (
	"errors"
	"net/http"
	"testing"
	
	"url-shortener-api/models"
)

func TestAppError_Error(t *testing.T) {
	appErr := &AppError{
		Message:    "test error message",
		StatusCode: http.StatusBadRequest,
	}

	expected := "test error message"
	if appErr.Error() != expected {
		t.Errorf("AppError.Error() = %v, want %v", appErr.Error(), expected)
	}
}

func TestAppError_GetStatusCode(t *testing.T) {
	appErr := &AppError{
		Message:    "test error message",
		StatusCode: http.StatusConflict,
	}

	expected := http.StatusConflict
	if appErr.GetStatusCode() != expected {
		t.Errorf("AppError.GetStatusCode() = %v, want %v", appErr.GetStatusCode(), expected)
	}
}

func TestGetStatusCodeFromError_AppError(t *testing.T) {
	tests := []struct {
		name           string
		error          error
		expectedStatus int
	}{
		{
			name:           "invalid URL format",
			error:          ErrInvalidURLFormat,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid URL scheme",
			error:          ErrInvalidURLScheme,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid alias length",
			error:          ErrInvalidAliasLength,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid alias characters",
			error:          ErrInvalidAliasChars,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "alias already exists",
			error:          ErrAliasAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "short code not found",
			error:          ErrShortCodeNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "short code expired",
			error:          ErrShortCodeExpired,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := GetStatusCodeFromError(tt.error)
			if statusCode != tt.expectedStatus {
				t.Errorf("GetStatusCodeFromError() = %d, want %d", statusCode, tt.expectedStatus)
			}
		})
	}
}

func TestGetStatusCodeFromError_UnknownError(t *testing.T) {
	unknownError := errors.New("unknown error")
	statusCode := GetStatusCodeFromError(unknownError)

	expected := http.StatusInternalServerError
	if statusCode != expected {
		t.Errorf("GetStatusCodeFromError() for unknown error = %d, want %d", statusCode, expected)
	}
}

func TestGetStatusCodeFromError_NilError(t *testing.T) {
	statusCode := GetStatusCodeFromError(nil)

	expected := http.StatusInternalServerError
	if statusCode != expected {
		t.Errorf("GetStatusCodeFromError() for nil error = %d, want %d", statusCode, expected)
	}
}

func TestErrorEquality(t *testing.T) {
	// Test that error instances are equal to themselves
	if ErrInvalidURLFormat != ErrInvalidURLFormat {
		t.Error("ErrInvalidURLFormat should equal itself")
	}

	if ErrAliasAlreadyExists != ErrAliasAlreadyExists {
		t.Error("ErrAliasAlreadyExists should equal itself")
	}

	// Test that different errors are not equal
	if ErrInvalidURLFormat == ErrAliasAlreadyExists {
		t.Error("Different errors should not be equal")
	}
}
