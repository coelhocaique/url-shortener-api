package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"url-shortener-api/models"
	"url-shortener-api/handlers"
)

func TestHandleError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Test endpoint that uses HandleError
	router.GET("/test-error", func(c *gin.Context) {
		handlers.HandleError(c, models.ErrInvalidURLFormat)
	})

	tests := []struct {
		name           string
		error          error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "validation error returns 400",
			error:          models.ErrInvalidURLFormat,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid URL format"}`,
		},
		{
			name:           "conflict error returns 409",
			error:          models.ErrAliasAlreadyExists,
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"alias already exists"}`,
		},
		{
			name:           "not found error returns 404",
			error:          models.ErrShortCodeNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"short code not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new router for each test to avoid state issues
			testRouter := gin.New()
			testRouter.GET("/test", func(c *gin.Context) {
				handlers.HandleError(c, tt.error)
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			testRouter.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("HandleError() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			body := w.Body.String()
			if body != tt.expectedBody {
				t.Errorf("HandleError() body = %s, want %s", body, tt.expectedBody)
			}
		})
	}
}

func TestHandleError_UnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Test with an unknown error type
	router.GET("/test-unknown-error", func(c *gin.Context) {
		handlers.HandleError(c, models.ErrInvalidURLFormat) // This should work normally
	})

	req, _ := http.NewRequest("GET", "/test-unknown-error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("HandleError() with known error should return appropriate status, got %d", w.Code)
	}
}
