package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"url-shortener-api/config"
	"url-shortener-api/routes"
	"url-shortener-api/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize services using factory
	serviceFactory := services.NewServiceFactory()
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
