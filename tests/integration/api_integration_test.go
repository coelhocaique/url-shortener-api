package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"url-shortener-api/models"
	"url-shortener-api/routes"
	"url-shortener-api/services"
)

func setupTestServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	// Initialize services
	serviceFactory := services.NewServiceFactory()
	urlService := serviceFactory.CreateURLService()
	
	// Setup router
	router := gin.Default()
	routes.SetupRoutes(router, urlService)
	
	return router
}

func TestAPIIntegration_CreateAndRedirectURL(t *testing.T) {
	router := setupTestServer()

	// Test 1: Create a short URL
	createRequest := models.URLRequest{
		URL:          "https://www.example.com",
		ExpirationMs: 3600000,
	}
	
	jsonBody, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify creation response
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response models.URLResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.ShortCode == "" {
		t.Error("Expected non-empty short code")
	}

	// Test 2: Redirect to the created URL
	redirectReq, _ := http.NewRequest("GET", fmt.Sprintf("/urls/%s", response.ShortCode), nil)
	redirectW := httptest.NewRecorder()
	router.ServeHTTP(redirectW, redirectReq)

	// Verify redirect response
	if redirectW.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status %d, got %d", http.StatusMovedPermanently, redirectW.Code)
	}

	location := redirectW.Header().Get("Location")
	if location != "https://www.example.com" {
		t.Errorf("Expected Location header 'https://www.example.com', got '%s'", location)
	}
}

func TestAPIIntegration_CreateURLWithCustomAlias(t *testing.T) {
	router := setupTestServer()

	// Test: Create URL with custom alias
	createRequest := models.URLRequest{
		URL:          "https://www.github.com",
		Alias:        "github",
		ExpirationMs: 7200000,
	}
	
	jsonBody, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify response
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response models.URLResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.ShortCode != "github" {
		t.Errorf("Expected short code 'github', got '%s'", response.ShortCode)
	}

	// Test redirect
	redirectReq, _ := http.NewRequest("GET", "/urls/github", nil)
	redirectW := httptest.NewRecorder()
	router.ServeHTTP(redirectW, redirectReq)

	if redirectW.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status %d, got %d", http.StatusMovedPermanently, redirectW.Code)
	}

	location := redirectW.Header().Get("Location")
	if location != "https://www.github.com" {
		t.Errorf("Expected Location header 'https://www.github.com', got '%s'", location)
	}
}

func TestAPIIntegration_ErrorScenarios(t *testing.T) {
	router := setupTestServer()

	tests := []struct {
		name           string
		requestBody    models.URLRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "invalid URL format",
			requestBody: models.URLRequest{
				URL:          "not a valid url",
				ExpirationMs: 3600000,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid URL format",
		},
		{
			name: "invalid alias length",
			requestBody: models.URLRequest{
				URL:          "https://www.example.com",
				Alias:        "ab",
				ExpirationMs: 3600000,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "alias must be between 3 and 20 characters",
		},
		{
			name: "invalid alias characters",
			requestBody: models.URLRequest{
				URL:          "https://www.example.com",
				Alias:        "invalid@alias",
				ExpirationMs: 3600000,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "alias can only contain letters, numbers, and hyphens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/urls", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var errorResponse map[string]string
			json.Unmarshal(w.Body.Bytes(), &errorResponse)
			if errorResponse["error"] != tt.expectedError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedError, errorResponse["error"])
			}
		})
	}
}

func TestAPIIntegration_DuplicateAlias(t *testing.T) {
	router := setupTestServer()

	// Create first URL with alias
	createRequest1 := models.URLRequest{
		URL:          "https://www.example1.com",
		Alias:        "duplicate-test",
		ExpirationMs: 3600000,
	}
	
	jsonBody1, _ := json.Marshal(createRequest1)
	req1, _ := http.NewRequest("POST", "/urls", bytes.NewBuffer(jsonBody1))
	req1.Header.Set("Content-Type", "application/json")
	
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Errorf("First creation should succeed, got status %d", w1.Code)
	}

	// Try to create second URL with same alias
	createRequest2 := models.URLRequest{
		URL:          "https://www.example2.com",
		Alias:        "duplicate-test",
		ExpirationMs: 3600000,
	}
	
	jsonBody2, _ := json.Marshal(createRequest2)
	req2, _ := http.NewRequest("POST", "/urls", bytes.NewBuffer(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")
	
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should get conflict error
	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w2.Code)
	}

	var errorResponse map[string]string
	json.Unmarshal(w2.Body.Bytes(), &errorResponse)
	if errorResponse["error"] != "alias already exists" {
		t.Errorf("Expected error 'alias already exists', got '%s'", errorResponse["error"])
	}
}

func TestAPIIntegration_ExpiredURL(t *testing.T) {
	router := setupTestServer()

	// Create URL with very short expiration
	createRequest := models.URLRequest{
		URL:          "https://www.example.com",
		Alias:        "expire-test",
		ExpirationMs: 1, // 1 millisecond
	}
	
	jsonBody, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Try to access expired URL
	redirectReq, _ := http.NewRequest("GET", "/urls/expire-test", nil)
	redirectW := httptest.NewRecorder()
	router.ServeHTTP(redirectW, redirectReq)

	// Should get not found error
	if redirectW.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, redirectW.Code)
	}

	var errorResponse map[string]string
	json.Unmarshal(redirectW.Body.Bytes(), &errorResponse)
	if errorResponse["error"] != "short code has expired" {
		t.Errorf("Expected error 'short code has expired', got '%s'", errorResponse["error"])
	}
}

func TestAPIIntegration_NonExistentShortCode(t *testing.T) {
	router := setupTestServer()

	// Try to access non-existent short code
	req, _ := http.NewRequest("GET", "/urls/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should get not found error
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var errorResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if errorResponse["error"] != "short code not found" {
		t.Errorf("Expected error 'short code not found', got '%s'", errorResponse["error"])
	}
}

func TestAPIIntegration_URLNormalization(t *testing.T) {
	router := setupTestServer()

	// Test URL without protocol
	createRequest := models.URLRequest{
		URL:          "www.example.com",
		ExpirationMs: 3600000,
	}
	
	jsonBody, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response models.URLResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.ShortCode == "" {
		t.Error("Expected non-empty short code")
	}

	// Test redirect to verify URL was normalized
	redirectReq, _ := http.NewRequest("GET", fmt.Sprintf("/urls/%s", response.ShortCode), nil)
	redirectW := httptest.NewRecorder()
	router.ServeHTTP(redirectW, redirectReq)

	if redirectW.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status %d, got %d", http.StatusMovedPermanently, redirectW.Code)
	}

	location := redirectW.Header().Get("Location")
	expected := "https://www.example.com"
	if location != expected {
		t.Errorf("Expected Location header '%s', got '%s'", expected, location)
	}
}
