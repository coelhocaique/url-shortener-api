package main

import (
	"net/http"
	"testing"
	"errors"
	"url-shortener-api/models"
)

func TestAppError_Error(t *testing.T) {
	appErr := &models.AppError{
		Message:    "test error message",
		StatusCode: http.StatusBadRequest,
	}

	expected := "test error message"
	if appErr.Error() != expected {
		t.Errorf("AppError.Error() = %v, want %v", appErr.Error(), expected)
	}
}

func TestAppError_GetStatusCode(t *testing.T) {
	appErr := &models.AppError{
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
			error:          models.ErrInvalidURLFormat,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid URL scheme",
			error:          models.ErrInvalidURLScheme,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid alias length",
			error:          models.ErrInvalidAliasLength,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid alias characters",
			error:          models.ErrInvalidAliasChars,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "alias already exists",
			error:          models.ErrAliasAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "short code not found",
			error:          models.ErrShortCodeNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "short code expired",
			error:          models.ErrShortCodeExpired,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := models.GetStatusCodeFromError(tt.error)
			if statusCode != tt.expectedStatus {
				t.Errorf("GetStatusCodeFromError() = %d, want %d", statusCode, tt.expectedStatus)
			}
		})
	}
}

func TestGetStatusCodeFromError_UnknownError(t *testing.T) {
	unknownError := errors.New("unknown error")
	statusCode := models.GetStatusCodeFromError(unknownError)

	expected := http.StatusInternalServerError
	if statusCode != expected {
		t.Errorf("GetStatusCodeFromError() for unknown error = %d, want %d", statusCode, expected)
	}
}

func TestGetStatusCodeFromError_NilError(t *testing.T) {
	statusCode := models.GetStatusCodeFromError(nil)

	expected := http.StatusInternalServerError
	if statusCode != expected {
		t.Errorf("GetStatusCodeFromError() for nil error = %d, want %d", statusCode, expected)
	}
}

func TestErrorEquality(t *testing.T) {
	// Test that error instances are equal to themselves
	if models.ErrInvalidURLFormat != models.ErrInvalidURLFormat {
		t.Error("ErrInvalidURLFormat should equal itself")
	}

	if models.ErrAliasAlreadyExists != models.ErrAliasAlreadyExists {
		t.Error("ErrAliasAlreadyExists should equal itself")
	}

	// Test that different errors are not equal
	if models.ErrInvalidURLFormat == models.ErrAliasAlreadyExists {
		t.Error("Different errors should not be equal")
	}
}
