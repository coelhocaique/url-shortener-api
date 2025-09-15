package testutils

import (
	"context"
	"testing"
	"time"

	"url-shortener-api/services"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/redis"
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

// SetupTestMongoDB creates a test MongoDB connection and collection using Docker
func SetupTestMongoDB(t *testing.T, config *TestMongoDBConfig) (*mongo.Client, *mongo.Collection, func()) {
	if config == nil {
		config = DefaultTestConfig()
	}

	ctx := context.Background()

	// Start MongoDB container
	mongoContainer, err := mongodb.RunContainer(ctx,
		testcontainers.WithImage("mongo:7.0"),
		mongodb.WithUsername("testuser"),
		mongodb.WithPassword("testpass"),
	)
	if err != nil {
		t.Fatalf("Failed to start MongoDB container: %v", err)
	}

	// Get connection string from container
	connStr, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get MongoDB connection string: %v", err)
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
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

		// Terminate the container
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate MongoDB container: %v", err)
		}
	}

	return client, collection, cleanup
}

// SetupTestRedis creates a test Redis connection using Docker
func SetupTestRedis(t *testing.T) (string, func()) {
	ctx := context.Background()

	// Start Redis container
	redisContainer, err := redis.RunContainer(ctx,
		testcontainers.WithImage("redis:7.2-alpine"),
	)
	if err != nil {
		t.Fatalf("Failed to start Redis container: %v", err)
	}

	// Get connection string from container
	connStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get Redis connection string: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		// Terminate the container
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate Redis container: %v", err)
		}
	}

	return connStr, cleanup
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
	_, collection, mongoCleanup := SetupTestMongoDB(t, nil)
	redisURL, redisCleanup := SetupTestRedis(t)

	factory := services.NewServiceFactory(collection, redisURL)

	// Combined cleanup function
	cleanup := func() {
		mongoCleanup()
		redisCleanup()
	}

	return factory, cleanup
}
