package services_test

import (
	"testing"
	"time"

	"url-shortener-api/models"
	"url-shortener-api/tests/testutils"
)

func TestURLStorage_StoreAndGet(t *testing.T) {
	storage, cleanup := testutils.CreateTestURLStorage(t)
	defer cleanup()

	expirationTime := time.Now().Add(time.Hour)
	mapping := models.URLMapping{
		OriginalURL:         "https://www.example.com",
		Alias:               "test",
		ExpirationTimestamp: &expirationTime,
		UserID:              "user123",
	}

	// Test storing
	err := storage.Store("test123", mapping)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	// Test getting
	retrieved, exists, err := storage.Get("test123")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !exists {
		t.Errorf("Get() should return true for existing key")
	}

	if retrieved.OriginalURL != mapping.OriginalURL {
		t.Errorf("Get() OriginalURL = %v, want %v", retrieved.OriginalURL, mapping.OriginalURL)
	}

	if retrieved.Alias != mapping.Alias {
		t.Errorf("Get() Alias = %v, want %v", retrieved.Alias, mapping.Alias)
	}

	if retrieved.UserID != mapping.UserID {
		t.Errorf("Get() UserID = %v, want %v", retrieved.UserID, mapping.UserID)
	}
}

func TestURLStorage_Exists(t *testing.T) {
	storage, cleanup := testutils.CreateTestURLStorage(t)
	defer cleanup()

	expirationTime := time.Now().Add(time.Hour)
	mapping := models.URLMapping{
		OriginalURL:         "https://www.example.com",
		ExpirationTimestamp: &expirationTime,
		UserID:              "user123",
	}

	// Test non-existent key
	exists, err := storage.Exists("nonexistent")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Errorf("Exists() should return false for non-existent key")
	}

	// Test existing key
	err = storage.Store("test", mapping)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	exists, err = storage.Exists("test")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Errorf("Exists() should return true for existing key")
	}
}

func TestURLStorage_Delete(t *testing.T) {
	storage, cleanup := testutils.CreateTestURLStorage(t)
	defer cleanup()

	expirationTime := time.Now().Add(time.Hour)
	mapping := models.URLMapping{
		OriginalURL:         "https://www.example.com",
		ExpirationTimestamp: &expirationTime,
		UserID:              "user123",
	}

	err := storage.Store("test", mapping)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	// Verify it exists
	exists, err := storage.Exists("test")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Errorf("Key should exist before deletion")
	}

	// Delete it
	err = storage.Delete("test")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify it no longer exists
	exists, err = storage.Exists("test")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Errorf("Key should not exist after deletion")
	}
}

func TestURLStorage_IsExpired(t *testing.T) {
	storage, cleanup := testutils.CreateTestURLStorage(t)
	defer cleanup()

	tests := []struct {
		name           string
		expirationTime *time.Time
		expected       bool
	}{
		{
			name:           "not expired",
			expirationTime: func() *time.Time { t := time.Now().Add(time.Hour); return &t }(),
			expected:       false,
		},
		{
			name:           "expired",
			expirationTime: func() *time.Time { t := time.Now().Add(-time.Hour); return &t }(),
			expected:       true,
		},
		{
			name:           "nil expiration (no expiration)",
			expirationTime: nil,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping := models.URLMapping{
				OriginalURL:         "https://www.example.com",
				ExpirationTimestamp: tt.expirationTime,
				UserID:              "user123",
			}

			result := storage.IsExpired(mapping)
			if result != tt.expected {
				t.Errorf("IsExpired() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestURLStorage_Overwrite(t *testing.T) {
	storage, cleanup := testutils.CreateTestURLStorage(t)
	defer cleanup()

	expirationTime1 := time.Now().Add(time.Hour)
	expirationTime2 := time.Now().Add(2 * time.Hour)

	originalMapping := models.URLMapping{
		OriginalURL:         "https://www.example.com",
		ExpirationTimestamp: &expirationTime1,
		UserID:              "user123",
	}

	newMapping := models.URLMapping{
		OriginalURL:         "https://www.newexample.com",
		ExpirationTimestamp: &expirationTime2,
		UserID:              "user456",
	}

	// Store original
	err := storage.Store("test", originalMapping)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	// Overwrite with new
	err = storage.Store("test", newMapping)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	// Verify new mapping is stored
	retrieved, exists, err := storage.Get("test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !exists {
		t.Errorf("Key should exist after overwrite")
	}

	if retrieved.OriginalURL != newMapping.OriginalURL {
		t.Errorf("Get() should return new OriginalURL after overwrite")
	}
}

// Test new MongoDB-specific methods
func TestURLStorage_GetByAlias(t *testing.T) {
	storage, cleanup := testutils.CreateTestURLStorage(t)
	defer cleanup()

	expirationTime := time.Now().Add(time.Hour)
	mapping := models.URLMapping{
		OriginalURL:         "https://www.example.com",
		Alias:               "test-alias",
		ExpirationTimestamp: &expirationTime,
		UserID:              "user123",
	}

	// Store with alias
	err := storage.Store("test123", mapping)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	// Get by alias
	retrieved, exists, err := storage.GetByAlias("test-alias")
	if err != nil {
		t.Fatalf("GetByAlias() error = %v", err)
	}
	if !exists {
		t.Errorf("GetByAlias() should return true for existing alias")
	}

	if retrieved.OriginalURL != mapping.OriginalURL {
		t.Errorf("GetByAlias() OriginalURL = %v, want %v", retrieved.OriginalURL, mapping.OriginalURL)
	}
}

func TestURLStorage_GetByUserID(t *testing.T) {
	storage, cleanup := testutils.CreateTestURLStorage(t)
	defer cleanup()

	expirationTime := time.Now().Add(time.Hour)

	// Store multiple URLs for same user
	mapping1 := models.URLMapping{
		OriginalURL:         "https://www.example1.com",
		ExpirationTimestamp: &expirationTime,
		UserID:              "user123",
	}

	mapping2 := models.URLMapping{
		OriginalURL:         "https://www.example2.com",
		ExpirationTimestamp: &expirationTime,
		UserID:              "user123",
	}

	err := storage.Store("test1", mapping1)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	err = storage.Store("test2", mapping2)
	if err != nil {
		t.Fatalf("Store() error = %v", err)
	}

	// Get by user ID
	mappings, err := storage.GetByUserID("user123")
	if err != nil {
		t.Fatalf("GetByUserID() error = %v", err)
	}

	if len(mappings) != 2 {
		t.Errorf("GetByUserID() returned %d mappings, want 2", len(mappings))
	}
}
