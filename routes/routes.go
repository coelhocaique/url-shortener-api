package routes

import (
	"url-shortener-api/handlers"
	"url-shortener-api/middleware"
	"url-shortener-api/models"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine, urlService models.URLService) {
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "url-shortener-api",
		})
	})

	// Create handlers
	urlHandler := handlers.NewURLHandler(urlService)

	// URL routes (authentication required)
	urls := r.Group("/urls")
	urls.Use(middleware.AuthMiddleware())
	{
		urls.POST("", urlHandler.CreateShortURL)
		urls.GET("/:short_code", urlHandler.RedirectToURL)
	}
}
