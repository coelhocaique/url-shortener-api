package handlers

import (
	"net/http"
	"strconv"

	"url-shortener-api/models"

	"github.com/gin-gonic/gin"
)

// URLHandler handles HTTP requests for URL operations
type URLHandler struct {
	urlService models.URLService
}

// NewURLHandler creates a new instance of URLHandler
func NewURLHandler(urlService models.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

// CreateShortURL handles POST /urls
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req models.URLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	response, err := h.urlService.CreateShortURL(&req, userIDStr)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// RedirectToURL handles GET /urls/{short_code}
func (h *URLHandler) RedirectToURL(c *gin.Context) {
	shortCode := c.Param("short_code")

	// Parse use_cache query parameter (default: true)
	useCache := true
	if useCacheStr := c.Query("use_cache"); useCacheStr != "" {
		if parsed, err := strconv.ParseBool(useCacheStr); err == nil {
			useCache = parsed
		}
	}

	originalURL, err := h.urlService.GetOriginalURL(shortCode, useCache)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Redirect to original URL
	c.Header("Location", originalURL)
	c.Status(http.StatusMovedPermanently)
}
