package services_test

import (
	"testing"
	"time"

	"url-shortener-api/models"
	"url-shortener-api/services"
)

func TestURLStorage_StoreAndGet(t *testing.T) {
	storage := services.NewURLStorage()
	
	mapping := models.URLMapping{
		OriginalURL:    "https://www.example.com",
		Alias:          "test",
		ExpirationTime: time.Now().Add(time.Hour),
	}

	// Test storing
	storage.Store("test123", mapping)

	// Test getting
	retrieved, exists := storage.Get("test123")
	if !exists {
		t.Errorf("Get() should return true for existing key")
	}

	if retrieved.OriginalURL != mapping.OriginalURL {
		t.Errorf("Get() OriginalURL = %v, want %v", retrieved.OriginalURL, mapping.OriginalURL)
	}

	if retrieved.Alias != mapping.Alias {
		t.Errorf("Get() Alias = %v, want %v", retrieved.Alias, mapping.Alias)
	}
}

func TestURLStorage_Exists(t *testing.T) {
	storage := services.NewURLStorage()
	
	mapping := models.URLMapping{
		OriginalURL:    "https://www.example.com",
		ExpirationTime: time.Now().Add(time.Hour),
	}

	// Test non-existent key
	if storage.Exists("nonexistent") {
		t.Errorf("Exists() should return false for non-existent key")
	}

	// Test existing key
	storage.Store("test", mapping)
	if !storage.Exists("test") {
		t.Errorf("Exists() should return true for existing key")
	}
}

func TestURLStorage_Delete(t *testing.T) {
	storage := services.NewURLStorage()
	
	mapping := models.URLMapping{
		OriginalURL:    "https://www.example.com",
		ExpirationTime: time.Now().Add(time.Hour),
	}

	storage.Store("test", mapping)
	
	// Verify it exists
	if !storage.Exists("test") {
		t.Errorf("Key should exist before deletion")
	}

	// Delete it
	storage.Delete("test")

	// Verify it no longer exists
	if storage.Exists("test") {
		t.Errorf("Key should not exist after deletion")
	}
}

func TestURLStorage_IsExpired(t *testing.T) {
	storage := services.NewURLStorage()
	
	tests := []struct {
		name           string
		expirationTime time.Time
		expected       bool
	}{
		{
			name:           "not expired",
			expirationTime: time.Now().Add(time.Hour),
			expected:       false,
		},
		{
			name:           "expired",
			expirationTime: time.Now().Add(-time.Hour),
			expected:       true,
		},
		{
			name:           "zero time (no expiration)",
			expirationTime: time.Time{},
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping := models.URLMapping{
				OriginalURL:    "https://www.example.com",
				ExpirationTime: tt.expirationTime,
			}

			result := storage.IsExpired(mapping)
			if result != tt.expected {
				t.Errorf("IsExpired() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestURLStorage_Overwrite(t *testing.T) {
	storage := services.NewURLStorage()
	
	originalMapping := models.URLMapping{
		OriginalURL:    "https://www.example.com",
		ExpirationTime: time.Now().Add(time.Hour),
	}

	newMapping := models.URLMapping{
		OriginalURL:    "https://www.newexample.com",
		ExpirationTime: time.Now().Add(2 * time.Hour),
	}

	// Store original
	storage.Store("test", originalMapping)
	
	// Overwrite with new
	storage.Store("test", newMapping)

	// Verify new mapping is stored
	retrieved, exists := storage.Get("test")
	if !exists {
		t.Errorf("Key should exist after overwrite")
	}

	if retrieved.OriginalURL != newMapping.OriginalURL {
		t.Errorf("Get() should return new OriginalURL after overwrite")
	}
}
