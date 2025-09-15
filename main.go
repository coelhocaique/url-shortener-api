package main

import (
	"context"
	"fmt"
	"log"

	"url-shortener-api/config"
	"url-shortener-api/routes"
	"url-shortener-api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	client, err := connectToMongoDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal("Failed to disconnect from MongoDB:", err)
		}
	}()

	// Get database and collection
	db := client.Database(cfg.DatabaseName)
	collection := db.Collection("url_mappings")

	// Initialize services using factory with MongoDB collection and Redis URL
	serviceFactory := services.NewServiceFactory(collection, cfg.RedisURL)
	urlService := serviceFactory.CreateURLService()

	// Setup Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, urlService)

	// Start server
	fmt.Printf("URL Shortener API starting on :%s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// connectToMongoDB establishes a connection to MongoDB
func connectToMongoDB(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Connected to MongoDB at %s\n", cfg.MongoURI)
	return client, nil
}
