package services

import (
	"context"
	"fmt"
	"time"

	"url-shortener-api/models"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DistributedCounter handles distributed counter operations with Redis and MongoDB
type DistributedCounter struct {
	redisClient *redis.Client
	collection  *mongo.Collection
	streamName  string
}

// NewDistributedCounter creates a new instance of DistributedCounter
func NewDistributedCounter(redisClient *redis.Client, collection *mongo.Collection) *DistributedCounter {
	return &DistributedCounter{
		redisClient: redisClient,
		collection:  collection,
		streamName:  "counter_replication",
	}
}

// GetNextCounter increments and returns the next counter value
func (dc *DistributedCounter) GetNextCounter() (int64, error) {
	ctx := context.Background()

	// Increment counter in Redis
	counter, err := dc.redisClient.Incr(ctx, "short_code_counter").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment counter in Redis: %w", err)
	}

	// Send counter update to Redis stream for async replication
	streamData := map[string]interface{}{
		"counter":   counter,
		"timestamp": time.Now().Unix(),
	}

	_, err = dc.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: dc.streamName,
		Values: streamData,
	}).Result()

	if err != nil {
		// Log error but don't fail the operation
		// In production, you might want to use a proper logger
		fmt.Printf("Warning: Failed to add to replication stream: %v\n", err)
	}

	return counter, nil
}

// GetCurrentCounter returns the current counter value from Redis
func (dc *DistributedCounter) GetCurrentCounter() (int64, error) {
	ctx := context.Background()

	counter, err := dc.redisClient.Get(ctx, "short_code_counter").Int64()
	if err == redis.Nil {
		// Counter doesn't exist in Redis, initialize from MongoDB
		return dc.initializeFromMongoDB()
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get counter from Redis: %w", err)
	}

	return counter, nil
}

// InitializeCounter initializes the counter in both Redis and MongoDB
func (dc *DistributedCounter) InitializeCounter() error {
	ctx := context.Background()

	// Check if counter already exists in MongoDB
	var counterDoc models.ShortCodeCounter
	err := dc.collection.FindOne(ctx, bson.M{}).Decode(&counterDoc)

	if err == mongo.ErrNoDocuments {
		// Initialize with 0
		counterDoc = models.ShortCodeCounter{
			Counter:   0,
			UpdatedAt: time.Now(),
		}

		_, err = dc.collection.InsertOne(ctx, counterDoc)
		if err != nil {
			return fmt.Errorf("failed to initialize counter in MongoDB: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check counter in MongoDB: %w", err)
	}

	// Set counter in Redis
	err = dc.redisClient.Set(ctx, "short_code_counter", counterDoc.Counter, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set counter in Redis: %w", err)
	}

	return nil
}

// initializeFromMongoDB initializes Redis counter from MongoDB
func (dc *DistributedCounter) initializeFromMongoDB() (int64, error) {
	ctx := context.Background()

	var counterDoc models.ShortCodeCounter
	err := dc.collection.FindOne(ctx, bson.M{}).Decode(&counterDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Initialize counter
			if initErr := dc.InitializeCounter(); initErr != nil {
				return 0, initErr
			}
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get counter from MongoDB: %w", err)
	}

	// Set counter in Redis
	err = dc.redisClient.Set(ctx, "short_code_counter", counterDoc.Counter, 0).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to set counter in Redis: %w", err)
	}

	return counterDoc.Counter, nil
}

// ReplicateToMongoDB replicates counter updates from Redis stream to MongoDB
func (dc *DistributedCounter) ReplicateToMongoDB() error {
	fmt.Println("Replicating counter to MongoDB")
	ctx := context.Background()

	// Read from stream
	streams, err := dc.redisClient.XRead(ctx, &redis.XReadArgs{
		Streams: []string{dc.streamName, "0"},
		Count:   100,
		Block:   0,
	}).Result()

	if err != nil {
		return fmt.Errorf("failed to read from stream: %w", err)
	}

	for _, stream := range streams {
		for _, message := range stream.Messages {
			counterStr, exists := message.Values["counter"]
			if !exists {
				continue
			}

			counter, ok := counterStr.(int64)
			if !ok {
				// Try to convert from string
				if counterStrStr, ok := counterStr.(string); ok {
					if parsedCounter, err := fmt.Sscanf(counterStrStr, "%d", &counter); err != nil || parsedCounter != 1 {
						continue
					}
				} else {
					continue
				}
			}

			// Update MongoDB
			_, err = dc.collection.UpdateOne(
				ctx,
				bson.M{},
				bson.M{
					"$set": bson.M{
						"counter":    counter,
						"updated_at": time.Now(),
					},
				},
				options.Update().SetUpsert(true),
			)

			if err != nil {
				fmt.Printf("Warning: Failed to replicate counter %d to MongoDB: %v\n", counter, err)
				continue
			}

			// Acknowledge the message
			dc.redisClient.XAck(ctx, dc.streamName, "replication_group", message.ID)
		}
	}

	return nil
}
