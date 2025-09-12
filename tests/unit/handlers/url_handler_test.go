package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"url-shortener-api/handlers"
	"url-shortener-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// MockURLService is a mock implementation of URLService
type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) CreateShortURL(req *models.URLRequest) (*models.URLResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.URLResponse), args.Error(1)
}

func (m *MockURLService) GetOriginalURL(shortCode string) (string, error) {
	args := m.Called(shortCode)
	return args.String(0), args.Error(1)
}

func (m *MockURLService) DeleteExpiredURL(shortCode string) {
	m.Called(shortCode)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestURLHandler_CreateShortURL_Success(t *testing.T) {
	// Setup
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)
	router := setupTestRouter()
	router.POST("/urls", handler.CreateShortURL)

	// Mock service response
	expectedResponse := &models.URLResponse{ShortCode: "abc123"}
	mockService.On("CreateShortURL", mock.AnythingOfType("*models.URLRequest")).Return(expectedResponse, nil)

	// Create request
	requestBody := models.URLRequest{
		URL:          "https://www.example.com",
		ExpirationMs: 3600000,
		UserID:       "user123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Make request
	req, _ := http.NewRequest("POST", "/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response models.URLResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.ShortCode != "abc123" {
		t.Errorf("Expected short code 'abc123', got '%s'", response.ShortCode)
	}

	mockService.AssertExpectations(t)
}

func TestURLHandler_CreateShortURL_ValidationError(t *testing.T) {
	// Setup
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)
	router := setupTestRouter()
	router.POST("/urls", handler.CreateShortURL)

	// Mock service error
	mockService.On("CreateShortURL", mock.AnythingOfType("*models.URLRequest")).Return(nil, models.ErrInvalidURLFormat)

	// Create request
	requestBody := models.URLRequest{
		URL:          "invalid url",
		ExpirationMs: 3600000,
		UserID:       "user123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Make request
	req, _ := http.NewRequest("POST", "/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var errorResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if errorResponse["error"] != "invalid URL format" {
		t.Errorf("Expected error 'invalid URL format', got '%s'", errorResponse["error"])
	}

	mockService.AssertExpectations(t)
}

func TestURLHandler_CreateShortURL_AliasConflict(t *testing.T) {
	// Setup
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)
	router := setupTestRouter()
	router.POST("/urls", handler.CreateShortURL)

	// Mock service error
	mockService.On("CreateShortURL", mock.AnythingOfType("*models.URLRequest")).Return(nil, models.ErrAliasAlreadyExists)

	// Create request
	requestBody := models.URLRequest{
		URL:          "https://www.example.com",
		Alias:        "existing-alias",
		ExpirationMs: 3600000,
		UserID:       "user123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Make request
	req, _ := http.NewRequest("POST", "/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
	}

	var errorResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if errorResponse["error"] != "alias already exists" {
		t.Errorf("Expected error 'alias already exists', got '%s'", errorResponse["error"])
	}

	mockService.AssertExpectations(t)
}

func TestURLHandler_RedirectToURL_Success(t *testing.T) {
	// Setup
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)
	router := setupTestRouter()
	router.GET("/urls/:short_code", handler.RedirectToURL)

	// Mock service response
	mockService.On("GetOriginalURL", "abc123").Return("https://www.example.com", nil)

	// Make request
	req, _ := http.NewRequest("GET", "/urls/abc123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status %d, got %d", http.StatusMovedPermanently, w.Code)
	}

	location := w.Header().Get("Location")
	if location != "https://www.example.com" {
		t.Errorf("Expected Location header 'https://www.example.com', got '%s'", location)
	}

	mockService.AssertExpectations(t)
}

func TestURLHandler_RedirectToURL_NotFound(t *testing.T) {
	// Setup
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)
	router := setupTestRouter()
	router.GET("/urls/:short_code", handler.RedirectToURL)

	// Mock service error
	mockService.On("GetOriginalURL", "nonexistent").Return("", models.ErrShortCodeNotFound)

	// Make request
	req, _ := http.NewRequest("GET", "/urls/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var errorResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if errorResponse["error"] != "short code not found" {
		t.Errorf("Expected error 'short code not found', got '%s'", errorResponse["error"])
	}

	mockService.AssertExpectations(t)
}

func TestURLHandler_RedirectToURL_Expired(t *testing.T) {
	// Setup
	mockService := new(MockURLService)
	handler := handlers.NewURLHandler(mockService)
	router := setupTestRouter()
	router.GET("/urls/:short_code", handler.RedirectToURL)

	// Mock service error
	mockService.On("GetOriginalURL", "expired").Return("", models.ErrShortCodeExpired)

	// Make request
	req, _ := http.NewRequest("GET", "/urls/expired", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var errorResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if errorResponse["error"] != "short code has expired" {
		t.Errorf("Expected error 'short code has expired', got '%s'", errorResponse["error"])
	}

	mockService.AssertExpectations(t)
}
