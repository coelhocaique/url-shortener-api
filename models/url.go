package models

import "time"

// URLRequest represents the request body for creating a short URL
type URLRequest struct {
	URL          string `json:"url" binding:"required"`
	Alias        string `json:"alias"`
	ExpirationMs int64  `json:"expiration_ms"`
}

// URLResponse represents the response for creating a short URL
type URLResponse struct {
	ShortCode string `json:"short_code"`
}

// URLMapping represents a URL mapping in memory
type URLMapping struct {
	OriginalURL    string
	Alias          string
	ExpirationTime time.Time
}

// URLService interface defines the contract for URL operations
type URLService interface {
	CreateShortURL(req *URLRequest) (*URLResponse, error)
	GetOriginalURL(shortCode string) (string, error)
	DeleteExpiredURL(shortCode string)
}
