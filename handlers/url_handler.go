package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"url-shortener-api/models"
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

	response, err := h.urlService.CreateShortURL(&req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// RedirectToURL handles GET /urls/{short_code}
func (h *URLHandler) RedirectToURL(c *gin.Context) {
	shortCode := c.Param("short_code")

	originalURL, err := h.urlService.GetOriginalURL(shortCode)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Redirect to original URL
	c.Header("Location", originalURL)
	c.Status(http.StatusMovedPermanently)
}
