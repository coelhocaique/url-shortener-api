package testutils

import (
	"context"
	"testing"
	"time"

	"url-shortener-api/services"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestMongoDBConfig holds configuration for test MongoDB connection
type TestMongoDBConfig struct {
	URI        string
	Database   string
	Collection string
}

// DefaultTestConfig returns default test configuration
func DefaultTestConfig() *TestMongoDBConfig {
	return &TestMongoDBConfig{
		URI:        "mongodb://localhost:27017",
		Database:   "url_shortener_test",
		Collection: "url_mappings_test",
	}
}

// SetupTestMongoDB creates a test MongoDB connection and collection
func SetupTestMongoDB(t *testing.T, config *TestMongoDBConfig) (*mongo.Client, *mongo.Collection, func()) {
	if config == nil {
		config = DefaultTestConfig()
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))
	if err != nil {
		t.Fatalf("Failed to connect to test MongoDB: %v", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to ping test MongoDB: %v", err)
	}

	// Get database and collection
	db := client.Database(config.Database)
	collection := db.Collection(config.Collection)

	// Cleanup function
	cleanup := func() {
		// Drop the test collection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		collection.Drop(ctx)

		// Disconnect from MongoDB
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		client.Disconnect(ctx)
	}

	return client, collection, cleanup
}

// CreateTestURLStorage creates a URLStorage instance for testing
func CreateTestURLStorage(t *testing.T) (*services.URLStorage, func()) {
	_, collection, cleanup := SetupTestMongoDB(t, nil)

	storage := services.NewURLStorage(collection)

	// Create indexes
	if err := storage.CreateIndexes(); err != nil {
		t.Fatalf("Failed to create indexes: %v", err)
	}

	return storage, cleanup
}

// CreateTestServiceFactory creates a service factory for testing
func CreateTestServiceFactory(t *testing.T) (*services.ServiceFactory, func()) {
	_, collection, cleanup := SetupTestMongoDB(t, nil)

	factory := services.NewServiceFactory(collection)

	return factory, cleanup
}
