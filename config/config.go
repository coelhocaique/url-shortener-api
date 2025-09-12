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

	timeout := 10 * time.Second

	return &Config{
		Port:         port,
		MongoURI:     mongoURI,
		DatabaseName: databaseName,
		Timeout:      timeout,
	}
}
