package services

import (
	"time"
	"url-shortener-api/models"
)

// URLStorage handles URL mapping storage operations
type URLStorage struct {
	urlMappings map[string]models.URLMapping
}

// NewURLStorage creates a new instance of URLStorage
func NewURLStorage() *URLStorage {
	return &URLStorage{
		urlMappings: make(map[string]models.URLMapping),
	}
}

// Store saves a URL mapping
func (s *URLStorage) Store(shortCode string, mapping models.URLMapping) {
	s.urlMappings[shortCode] = mapping
}

// Get retrieves a URL mapping by short code
func (s *URLStorage) Get(shortCode string) (models.URLMapping, bool) {
	mapping, exists := s.urlMappings[shortCode]
	return mapping, exists
}

// Exists checks if a short code already exists
func (s *URLStorage) Exists(shortCode string) bool {
	_, exists := s.urlMappings[shortCode]
	return exists
}

// Delete removes a URL mapping
func (s *URLStorage) Delete(shortCode string) {
	delete(s.urlMappings, shortCode)
}

// IsExpired checks if a URL mapping has expired
func (s *URLStorage) IsExpired(mapping models.URLMapping) bool {
	return !mapping.ExpirationTime.IsZero() && time.Now().After(mapping.ExpirationTime)
}
