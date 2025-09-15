package config

import (
	"os"
	"time"
)

// Config holds application configuration
type Config struct {
	Port         string
	MongoURI     string
	DatabaseName string
	RedisURL     string
	Timeout      time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		databaseName = "url_shortener"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	timeout := 10 * time.Second

	return &Config{
		Port:         port,
		MongoURI:     mongoURI,
		DatabaseName: databaseName,
		RedisURL:     redisURL,
		Timeout:      timeout,
	}
}
