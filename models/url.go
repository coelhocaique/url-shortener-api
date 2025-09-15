package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

// URLMapping represents a URL mapping document in MongoDB
type URLMapping struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OriginalURL         string             `bson:"original_url" json:"original_url"`
	ShortURL            string             `bson:"short_url" json:"short_url"`
	ExpirationTimestamp *time.Time         `bson:"expiration_timestamp,omitempty" json:"expiration_timestamp,omitempty"`
	Alias               string             `bson:"alias,omitempty" json:"alias,omitempty"`
	CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at" json:"updated_at"`
	UserID              string             `bson:"user_id" json:"user_id"`
}

// URLService interface defines the contract for URL operations
type URLService interface {
	CreateShortURL(req *URLRequest, userID string) (*URLResponse, error)
	GetOriginalURL(shortCode string) (string, error)
	DeleteExpiredURL(shortCode string)
}
