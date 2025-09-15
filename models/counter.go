package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShortCodeCounter represents the counter document in MongoDB
type ShortCodeCounter struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Counter   int64              `bson:"counter" json:"counter"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// CounterService interface defines the contract for counter operations
type CounterService interface {
	GetNextCounter() (int64, error)
	GetCurrentCounter() (int64, error)
	InitializeCounter() error
}
