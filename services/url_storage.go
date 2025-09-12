package services

import (
	"context"
	"time"

	"url-shortener-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// URLStorage handles URL mapping storage operations with MongoDB
type URLStorage struct {
	collection *mongo.Collection
}

// NewURLStorage creates a new instance of URLStorage with MongoDB collection
func NewURLStorage(collection *mongo.Collection) *URLStorage {
	return &URLStorage{
		collection: collection,
	}
}

// Store saves a URL mapping to MongoDB
func (s *URLStorage) Store(shortCode string, mapping models.URLMapping) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set timestamps
	now := time.Now()
	mapping.CreatedAt = now
	mapping.UpdatedAt = now
	mapping.ShortURL = shortCode

	// Insert the document
	_, err := s.collection.InsertOne(ctx, mapping)
	return err
}

// Get retrieves a URL mapping by short code from MongoDB
func (s *URLStorage) Get(shortCode string) (models.URLMapping, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var mapping models.URLMapping
	filter := bson.M{"short_url": shortCode}

	err := s.collection.FindOne(ctx, filter).Decode(&mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return mapping, false, nil
		}
		return mapping, false, err
	}

	return mapping, true, nil
}

// Exists checks if a short code already exists in MongoDB
func (s *URLStorage) Exists(shortCode string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"short_url": shortCode}
	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Delete removes a URL mapping from MongoDB
func (s *URLStorage) Delete(shortCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"short_url": shortCode}
	_, err := s.collection.DeleteOne(ctx, filter)
	return err
}

// IsExpired checks if a URL mapping has expired
func (s *URLStorage) IsExpired(mapping models.URLMapping) bool {
	if mapping.ExpirationTimestamp == nil {
		return false
	}
	return time.Now().After(*mapping.ExpirationTimestamp)
}

// GetByAlias retrieves a URL mapping by alias from MongoDB
func (s *URLStorage) GetByAlias(alias string) (models.URLMapping, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var mapping models.URLMapping
	filter := bson.M{"alias": alias}

	err := s.collection.FindOne(ctx, filter).Decode(&mapping)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return mapping, false, nil
		}
		return mapping, false, err
	}

	return mapping, true, nil
}

// GetByUserID retrieves all URL mappings for a specific user
func (s *URLStorage) GetByUserID(userID string) ([]models.URLMapping, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mappings []models.URLMapping
	if err = cursor.All(ctx, &mappings); err != nil {
		return nil, err
	}

	return mappings, nil
}

// Update updates an existing URL mapping
func (s *URLStorage) Update(shortCode string, mapping models.URLMapping) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mapping.UpdatedAt = time.Now()
	filter := bson.M{"short_url": shortCode}
	update := bson.M{"$set": mapping}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

// CreateIndexes creates necessary indexes for the collection
func (s *URLStorage) CreateIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create index on short_url for fast lookups
	shortURLIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "short_url", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	// Create index on alias for fast lookups
	aliasIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "alias", Value: 1}},
		Options: options.Index().SetUnique(true).SetSparse(true),
	}

	// Create index on user_id for user-specific queries
	userIDIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	}

	// Create TTL index on expiration_timestamp for automatic cleanup
	ttlIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "expiration_timestamp", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err := s.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		shortURLIndex,
		aliasIndex,
		userIDIndex,
		ttlIndex,
	})

	return err
}
