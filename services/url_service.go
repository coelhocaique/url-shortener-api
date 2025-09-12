package services

import (
	"time"

	"url-shortener-api/models"
)

// URLServiceImpl implements the URLService interface
type URLServiceImpl struct {
	storage   *URLStorage
	generator *ShortCodeGenerator
	validator *URLValidator
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
		exists, err := s.storage.Exists(req.Alias)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, models.ErrAliasAlreadyExists
		}
		shortCode = req.Alias
	} else {
		// Generate unique short code
		for {
			shortCode = s.generator.Generate()
			exists, err := s.storage.Exists(shortCode)
			if err != nil {
				return nil, err
			}
			if !exists {
				break
			}
		}
	}

	// Calculate expiration time
	var expirationTime *time.Time
	if req.ExpirationMs > 0 {
		exp := time.Now().Add(time.Duration(req.ExpirationMs) * time.Millisecond)
		expirationTime = &exp
	}

	// Store the mapping
	mapping := models.URLMapping{
		OriginalURL:         validatedURL,
		Alias:               req.Alias,
		ExpirationTimestamp: expirationTime,
		UserID:              req.UserID,
	}

	if err := s.storage.Store(shortCode, mapping); err != nil {
		return nil, err
	}

	// Return response
	return &models.URLResponse{
		ShortCode: shortCode,
	}, nil
}

// GetOriginalURL retrieves the original URL for a given short code
func (s *URLServiceImpl) GetOriginalURL(shortCode string) (string, error) {
	// Check if mapping exists
	mapping, exists, err := s.storage.Get(shortCode)
	if err != nil {
		return "", err
	}
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
	if err := s.storage.Delete(shortCode); err != nil {
		// Log error but don't return it as this is cleanup
		// In a production app, you might want to use a proper logger
	}
}
