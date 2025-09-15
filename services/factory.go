package services

import (
	"log"
	"url-shortener-api/models"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

// ServiceFactory handles the creation and dependency injection of services
type ServiceFactory struct {
	collection         *mongo.Collection
	cache              *CacheService
	counterCollection  *mongo.Collection
	redisClient        *redis.Client
	replicationService *ReplicationService
}

// NewServiceFactory creates a new instance of ServiceFactory with MongoDB collection and Redis cache
func NewServiceFactory(collection *mongo.Collection, redisURL string) *ServiceFactory {
	cache := NewCacheService(redisURL)

	// Parse Redis URL to get client
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		// Fallback to default localhost connection
		opt = &redis.Options{
			Addr: "localhost:6379",
		}
	}
	redisClient := redis.NewClient(opt)

	// Create counter collection (using same database as main collection)
	db := collection.Database()
	counterCollection := db.Collection("short_code_counters")

	// Create distributed counter
	distributedCounter := NewDistributedCounter(redisClient, counterCollection)

	// Initialize counter
	if err := distributedCounter.InitializeCounter(); err != nil {
		log.Printf("Warning: Failed to initialize counter: %v", err)
	}

	// Create and start replication service
	replicationService := NewReplicationService(distributedCounter)
	replicationService.Start()

	return &ServiceFactory{
		collection:         collection,
		cache:              cache,
		counterCollection:  counterCollection,
		redisClient:        redisClient,
		replicationService: replicationService,
	}
}

// CreateURLService creates a new URLService with all its dependencies
func (f *ServiceFactory) CreateURLService() models.URLService {
	storage := NewURLStorage(f.collection)

	// Create distributed counter
	distributedCounter := NewDistributedCounter(f.redisClient, f.counterCollection)
	generator := NewShortCodeGenerator(distributedCounter)
	validator := NewURLValidator()

	// Create indexes for the collection
	if err := storage.CreateIndexes(); err != nil {
		log.Printf("Warning: Failed to create MongoDB indexes: %v", err)
	}

	return &URLServiceImpl{
		storage:   storage,
		generator: generator,
		validator: validator,
		cache:     f.cache,
	}
}
