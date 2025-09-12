package services

import (
	"log"
	"url-shortener-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

// ServiceFactory handles the creation and dependency injection of services
type ServiceFactory struct {
	collection *mongo.Collection
}

// NewServiceFactory creates a new instance of ServiceFactory with MongoDB collection
func NewServiceFactory(collection *mongo.Collection) *ServiceFactory {
	return &ServiceFactory{
		collection: collection,
	}
}

// CreateURLService creates a new URLService with all its dependencies
func (f *ServiceFactory) CreateURLService() models.URLService {
	storage := NewURLStorage(f.collection)
	generator := NewShortCodeGenerator()
	validator := NewURLValidator()

	// Create indexes for the collection
	if err := storage.CreateIndexes(); err != nil {
		log.Printf("Warning: Failed to create MongoDB indexes: %v", err)
	}

	return &URLServiceImpl{
		storage:   storage,
		generator: generator,
		validator: validator,
	}
}
