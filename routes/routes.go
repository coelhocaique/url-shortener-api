package routes

import (
	"github.com/gin-gonic/gin"
	"url-shortener-api/handlers"
	"url-shortener-api/models"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine, urlService models.URLService) {
	// Create handlers
	urlHandler := handlers.NewURLHandler(urlService)

	// URL routes
	urls := r.Group("/urls")
	{
		urls.POST("", urlHandler.CreateShortURL)
		urls.GET("/:short_code", urlHandler.RedirectToURL)
	}
}
