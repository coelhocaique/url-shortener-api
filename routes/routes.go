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

	// URL creation route (authentication required)
	urls := r.Group("/urls")
	urls.Use(middleware.AuthMiddleware())
	{
		urls.POST("", urlHandler.CreateShortURL)
	}

	// URL redirect route (no authentication required)
	r.GET("/urls/:short_code", urlHandler.RedirectToURL)

}
