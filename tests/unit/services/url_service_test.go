package services_test

import (
	"testing"
	"time"

	"url-shortener-api/models"
	"url-shortener-api/tests/testutils"
)

func TestURLServiceImpl_CreateShortURL(t *testing.T) {
	factory, cleanup := testutils.CreateTestServiceFactory(t)
	defer cleanup()
	service := factory.CreateURLService()

	tests := []struct {
		name        string
		request     *models.URLRequest
		userID      string
		wantErr     bool
		expectedErr error
	}{
		{
			name: "valid request without alias",
			request: &models.URLRequest{
				URL:          "https://www.example.com",
				ExpirationMs: 3600000,
			},
			userID:  "user123",
			wantErr: false,
		},
		{
			name: "valid request with alias",
			request: &models.URLRequest{
				URL:          "https://www.example.com",
				Alias:        "test-alias",
				ExpirationMs: 3600000,
			},
			userID:  "user123",
			wantErr: false,
		},
		{
			name: "invalid URL format",
			request: &models.URLRequest{
				URL:          "not a url",
				ExpirationMs: 3600000,
			},
			userID:      "user123",
			wantErr:     true,
			expectedErr: models.ErrInvalidURLFormat,
		},
		{
			name: "invalid alias length",
			request: &models.URLRequest{
				URL:          "https://www.example.com",
				Alias:        "ab",
				ExpirationMs: 3600000,
			},
			userID:      "user123",
			wantErr:     true,
			expectedErr: models.ErrInvalidAliasLength,
		},
		{
			name: "invalid alias characters",
			request: &models.URLRequest{
				URL:          "https://www.example.com",
				Alias:        "invalid@alias",
				ExpirationMs: 3600000,
			},
			userID:      "user123",
			wantErr:     true,
			expectedErr: models.ErrInvalidAliasChars,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.CreateShortURL(tt.request, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateShortURL() expected error but got none")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("CreateShortURL() error = %v, want %v", err, tt.expectedErr)
				}
			} else {
				if err != nil {
					t.Errorf("CreateShortURL() unexpected error = %v", err)
					return
				}
				if response == nil {
					t.Errorf("CreateShortURL() response is nil")
					return
				}
				if response.ShortCode == "" {
					t.Errorf("CreateShortURL() short code is empty")
				}
			}
		})
	}
}

func TestURLServiceImpl_CreateShortURLWithAlias(t *testing.T) {
	factory, cleanup := testutils.CreateTestServiceFactory(t)
	defer cleanup()
	service := factory.CreateURLService()

	// Create first URL with alias
	request1 := &models.URLRequest{
		URL:          "https://www.example1.com",
		Alias:        "test-alias",
		ExpirationMs: 3600000,
	}

	response1, err := service.CreateShortURL(request1, "user123")
	if err != nil {
		t.Fatalf("Failed to create first URL: %v", err)
	}

	if response1.ShortCode != "test-alias" {
		t.Errorf("CreateShortURL() short code = %v, want %v", response1.ShortCode, "test-alias")
	}

	// Try to create second URL with same alias
	request2 := &models.URLRequest{
		URL:          "https://www.example2.com",
		Alias:        "test-alias",
		ExpirationMs: 3600000,
	}

	_, err = service.CreateShortURL(request2, "user456")
	if err != models.ErrAliasAlreadyExists {
		t.Errorf("CreateShortURL() error = %v, want %v", err, models.ErrAliasAlreadyExists)
	}
}

func TestURLServiceImpl_GetOriginalURL(t *testing.T) {
	factory, cleanup := testutils.CreateTestServiceFactory(t)
	defer cleanup()
	service := factory.CreateURLService()

	// Create a URL first
	request := &models.URLRequest{
		URL:          "https://www.example.com",
		Alias:        "test-get",
		ExpirationMs: 3600000,
	}

	response, err := service.CreateShortURL(request, "user123")
	if err != nil {
		t.Fatalf("Failed to create URL: %v", err)
	}

	// Test getting the URL
	originalURL, err := service.GetOriginalURL(response.ShortCode)
	if err != nil {
		t.Errorf("GetOriginalURL() error = %v", err)
	}

	if originalURL != "https://www.example.com" {
		t.Errorf("GetOriginalURL() = %v, want %v", originalURL, "https://www.example.com")
	}
}

func TestURLServiceImpl_GetOriginalURLNotFound(t *testing.T) {
	factory, cleanup := testutils.CreateTestServiceFactory(t)
	defer cleanup()
	service := factory.CreateURLService()

	_, err := service.GetOriginalURL("nonexistent")
	if err != models.ErrShortCodeNotFound {
		t.Errorf("GetOriginalURL() error = %v, want %v", err, models.ErrShortCodeNotFound)
	}
}

func TestURLServiceImpl_GetOriginalURLExpired(t *testing.T) {
	factory, cleanup := testutils.CreateTestServiceFactory(t)
	defer cleanup()
	service := factory.CreateURLService()

	// Create a URL with very short expiration
	request := &models.URLRequest{
		URL:          "https://www.example.com",
		Alias:        "expired-test",
		ExpirationMs: 1, // 1 millisecond
	}

	response, err := service.CreateShortURL(request, "user123")
	if err != nil {
		t.Fatalf("Failed to create URL: %v", err)
	}

	// Wait for it to expire
	time.Sleep(10 * time.Millisecond)

	// Try to get the expired URL
	_, err = service.GetOriginalURL(response.ShortCode)
	if err != models.ErrShortCodeExpired {
		t.Errorf("GetOriginalURL() error = %v, want %v", err, models.ErrShortCodeExpired)
	}
}

func TestURLServiceImpl_URLNormalization(t *testing.T) {
	factory, cleanup := testutils.CreateTestServiceFactory(t)
	defer cleanup()
	service := factory.CreateURLService()

	// Test URL without protocol
	request := &models.URLRequest{
		URL:          "www.example.com",
		ExpirationMs: 3600000,
	}

	response, err := service.CreateShortURL(request, "user123")
	if err != nil {
		t.Fatalf("Failed to create URL: %v", err)
	}

	// Get the URL and verify it was normalized
	originalURL, err := service.GetOriginalURL(response.ShortCode)
	if err != nil {
		t.Errorf("GetOriginalURL() error = %v", err)
	}

	expected := "https://www.example.com"
	if originalURL != expected {
		t.Errorf("GetOriginalURL() = %v, want %v", originalURL, expected)
	}
}
