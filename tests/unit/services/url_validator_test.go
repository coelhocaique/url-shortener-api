package services_test

import (
	"testing"

	"url-shortener-api/models"
	"url-shortener-api/services"
)

func TestURLValidator_ValidateURL(t *testing.T) {
	validator := services.NewURLValidator()

	tests := []struct {
		name     string
		url      string
		expected string
		wantErr  bool
		errType  error
	}{
		{
			name:     "valid https URL",
			url:      "https://www.example.com",
			expected: "https://www.example.com",
			wantErr:  false,
		},
		{
			name:     "valid http URL",
			url:      "http://www.example.com",
			expected: "http://www.example.com",
			wantErr:  false,
		},
		{
			name:     "URL without protocol",
			url:      "www.example.com",
			expected: "https://www.example.com",
			wantErr:  false,
		},
		{
			name:     "URL with path",
			url:      "example.com/path",
			expected: "https://example.com/path",
			wantErr:  false,
		},
		{
			name:     "URL with query parameters",
			url:      "example.com?param=value",
			expected: "https://example.com?param=value",
			wantErr:  false,
		},
		{
			name:    "invalid URL format",
			url:     "not a url",
			wantErr: true,
			errType: models.ErrInvalidURLFormat,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
			errType: models.ErrInvalidURLFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.ValidateURL(tt.url)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateURL() expected error but got none")
					return
				}
				if tt.errType != nil && err != tt.errType {
					t.Errorf("ValidateURL() error = %v, want %v", err, tt.errType)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateURL() unexpected error = %v", err)
					return
				}
				if result != tt.expected {
					t.Errorf("ValidateURL() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestURLValidator_ValidateAlias(t *testing.T) {
	validator := services.NewURLValidator()

	tests := []struct {
		name    string
		alias   string
		wantErr bool
		errType error
	}{
		{
			name:    "valid alias with letters and numbers",
			alias:   "my-alias123",
			wantErr: false,
		},
		{
			name:    "valid alias with hyphens",
			alias:   "my-alias-name",
			wantErr: false,
		},
		{
			name:    "empty alias (valid)",
			alias:   "",
			wantErr: false,
		},
		{
			name:    "minimum length alias",
			alias:   "abc",
			wantErr: false,
		},
		{
			name:    "maximum length alias",
			alias:   "abcdefghijklmnopqrst",
			wantErr: false,
		},
		{
			name:    "alias too short",
			alias:   "ab",
			wantErr: true,
			errType: models.ErrInvalidAliasLength,
		},
		{
			name:    "alias too long",
			alias:   "abcdefghijklmnopqrstuvwxyz",
			wantErr: true,
			errType: models.ErrInvalidAliasLength,
		},
		{
			name:    "alias with invalid characters",
			alias:   "invalid@alias",
			wantErr: true,
			errType: models.ErrInvalidAliasChars,
		},
		{
			name:    "alias with spaces",
			alias:   "invalid alias",
			wantErr: true,
			errType: models.ErrInvalidAliasChars,
		},
		{
			name:    "alias with underscore",
			alias:   "invalid_alias",
			wantErr: true,
			errType: models.ErrInvalidAliasChars,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAlias(tt.alias)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateAlias() expected error but got none")
					return
				}
				if tt.errType != nil && err != tt.errType {
					t.Errorf("ValidateAlias() error = %v, want %v", err, tt.errType)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateAlias() unexpected error = %v", err)
				}
			}
		})
	}
}
