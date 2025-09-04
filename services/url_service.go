package services

import (
	"time"

	"url-shortener-api/models"
)

// URLServiceImpl implements the URLService interface
type URLServiceImpl struct {
	storage    *URLStorage
	generator  *ShortCodeGenerator
	validator  *URLValidator
}

// NewURLService creates a new instance of URLServiceImpl
func NewURLService() *URLServiceImpl {
	return &URLServiceImpl{
		storage:    NewURLStorage(),
		generator:  NewShortCodeGenerator(),
		validator:  NewURLValidator(),
	}
}

// CreateShortURL creates a new short URL mapping
func (s *URLServiceImpl) CreateShortURL(req *models.URLRequest) (*models.URLResponse, error) {
	// Validate URL
	validatedURL, err := s.validator.ValidateURL(req.URL)
	if err != nil {
		return nil, err
	}

	// Validate alias if provided
	if err := s.validator.ValidateAlias(req.Alias); err != nil {
		return nil, err
	}

	// Generate short code
	var shortCode string
	if req.Alias != "" {
			// Check if alias already exists
	if s.storage.Exists(req.Alias) {
		return nil, models.ErrAliasAlreadyExists
	}
		shortCode = req.Alias
	} else {
		shortCode = s.generator.Generate()
	}

	// Calculate expiration time
	expirationTime := time.Now().Add(time.Duration(req.ExpirationMs) * time.Millisecond)

	// Store the mapping
	mapping := models.URLMapping{
		OriginalURL:    validatedURL,
		Alias:          req.Alias,
		ExpirationTime: expirationTime,
	}
	s.storage.Store(shortCode, mapping)

	// Return response
	return &models.URLResponse{
		ShortCode: shortCode,
	}, nil
}

// GetOriginalURL retrieves the original URL for a given short code
func (s *URLServiceImpl) GetOriginalURL(shortCode string) (string, error) {
	// Check if mapping exists
	mapping, exists := s.storage.Get(shortCode)
	if !exists {
		return "", models.ErrShortCodeNotFound
	}

	// Check if URL has expired
	if s.storage.IsExpired(mapping) {
		s.DeleteExpiredURL(shortCode)
		return "", models.ErrShortCodeExpired
	}

	return mapping.OriginalURL, nil
}

// DeleteExpiredURL removes an expired URL mapping
func (s *URLServiceImpl) DeleteExpiredURL(shortCode string) {
	s.storage.Delete(shortCode)
}
