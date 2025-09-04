package handlers

import (
	"github.com/gin-gonic/gin"
	"url-shortener-api/models"
)

// HandleError handles errors consistently across all handlers
func HandleError(c *gin.Context, err error) {
	statusCode := models.GetStatusCodeFromError(err)
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
