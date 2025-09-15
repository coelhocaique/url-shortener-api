package services

import (
	"context"
	"fmt"
	"time"

	"url-shortener-api/models"
)

// URLServiceImpl implements the URLService interface
type URLServiceImpl struct {
	storage   *URLStorage
	generator *ShortCodeGenerator
	validator *URLValidator
	cache     *CacheService
}

// CreateShortURL creates a new short URL mapping
func (s *URLServiceImpl) CreateShortURL(req *models.URLRequest, userID string) (*models.URLResponse, error) {
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
		// Generate unique short code using distributed counter
		generatedCode, err := s.generator.Generate()
		if err != nil {
			return nil, err
		}
		shortCode = generatedCode
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
		UserID:              userID,
	}

	if err := s.storage.Store(shortCode, mapping); err != nil {
		return nil, err
	}

	// Write through to cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("url:%s", shortCode)
	cacheValue := mapping.OriginalURL

	// Calculate TTL for cache
	var ttl time.Duration
	if expirationTime != nil {
		ttl = time.Until(*expirationTime)
		if ttl > 0 {
			// Only cache if not already expired
			s.cache.Set(ctx, cacheKey, cacheValue, ttl)
		}
	} else {
		// No expiration, cache for a long time (24 hours)
		s.cache.Set(ctx, cacheKey, cacheValue, 24*time.Hour)
	}

	// Return response
	return &models.URLResponse{
		ShortCode: shortCode,
	}, nil
}

// GetOriginalURL retrieves the original URL for a given short code
func (s *URLServiceImpl) GetOriginalURL(shortCode string, useCache bool) (string, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("url:%s", shortCode)

	// If cache is enabled, try to get from cache first
	if useCache {
		fmt.Println("Getting from cache")
		cachedURL, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			// Cache hit, return the URL
			return cachedURL, nil
		}
	}

	// Cache miss or cache disabled, fallback to MongoDB
	mapping, exists, err := s.storage.Get(shortCode)
	fmt.Println("Error getting from MongoDB")
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

	// If cache is enabled, cache the result for future requests
	if useCache {
		var ttl time.Duration
		if mapping.ExpirationTimestamp != nil {
			ttl = time.Until(*mapping.ExpirationTimestamp)
			if ttl > 0 {
				s.cache.Set(ctx, cacheKey, mapping.OriginalURL, ttl)
			}
		} else {
			// No expiration, cache for a long time (24 hours)
			s.cache.Set(ctx, cacheKey, mapping.OriginalURL, 24*time.Hour)
		}
	}

	return mapping.OriginalURL, nil
}

// DeleteExpiredURL removes an expired URL mapping
func (s *URLServiceImpl) DeleteExpiredURL(shortCode string) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("url:%s", shortCode)

	// Remove from cache
	s.cache.Delete(ctx, cacheKey)

	// Remove from storage
	if err := s.storage.Delete(shortCode); err != nil {
		// Log error but don't return it as this is cleanup
		// In a production app, you might want to use a proper logger
	}
}
