package services

import (
	"errors"
	"net/url"
	"strings"
)

// URLValidator handles URL validation operations
type URLValidator struct{}

// NewURLValidator creates a new instance of URLValidator
func NewURLValidator() *URLValidator {
	return &URLValidator{}
}

// ValidateURL checks if a URL is valid and formats it properly
func (v *URLValidator) ValidateURL(urlStr string) (string, error) {
	// Add protocol if missing
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	// Parse and validate URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", errors.New("invalid URL format")
	}

	// Check if URL has a valid scheme and host
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", errors.New("invalid URL: missing scheme or host")
	}

	return urlStr, nil
}

// ValidateAlias checks if an alias is valid
func (v *URLValidator) ValidateAlias(alias string) error {
	if alias == "" {
		return nil // Empty alias is valid (will generate random code)
	}

	// Check length
	if len(alias) < 3 || len(alias) > 20 {
		return errors.New("alias must be between 3 and 20 characters")
	}

	// Check for valid characters (alphanumeric and hyphens only)
	for _, char := range alias {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || 
			 char == '-') {
			return errors.New("alias can only contain letters, numbers, and hyphens")
		}
	}

	return nil
}
