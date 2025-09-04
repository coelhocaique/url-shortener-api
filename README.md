# URL Shortener API

A modular URL shortener API built with Go and Gin framework, following clean architecture principles.

## Features

- Create short URLs with optional custom aliases
- Set expiration times for URLs
- Automatic redirect to original URLs
- In-memory storage (for development)
- Modular architecture with separation of concerns

## Project Structure

```
url-shortener-api/
├── config/          # Configuration management
├── handlers/        # HTTP request handlers
│   ├── error_handler.go     # Error handling utilities
│   └── url_handler.go       # URL operation handlers
├── models/          # Data structures and interfaces
│   ├── errors.go            # Custom error types with HTTP status codes
│   └── url.go               # URL data structures
├── routes/          # Route definitions
├── services/        # Business logic
│   ├── factory.go           # Service dependency injection
│   ├── short_code_generator.go  # Short code generation
│   ├── url_service.go       # Main URL service
│   ├── url_storage.go       # URL storage operations
│   └── url_validator.go     # URL validation
├── tests/           # Test suite
│   ├── unit/                # Unit tests
│   │   ├── services/        # Service layer tests
│   │   ├── handlers/        # Handler layer tests
│   │   └── models/          # Model layer tests
│   └── integration/         # Integration tests
│       └── api_integration_test.go
├── main.go         # Application entry point
├── run_tests.sh    # Test runner script
└── README.md       # This file
```

## Setup

1. Make sure you have Go 1.21+ installed
2. Navigate to the project directory:
   ```bash
   cd url-shortener-api
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

The API will start on `http://localhost:8080` (or the port specified in the PORT environment variable)

## API Endpoints

### POST /urls

Create a new short URL.

**Request Body:**
```json
{
   "url": "www.google.com", 
   "alias": "my-alias",
   "expiration_ms": "1000000"
}
```

**Response (201 Created):**
```json
{
   "short_code": "P89g2" 
}
```

**Parameters:**
- `url` (required): The original URL to shorten
- `alias` (optional): Custom short code/alias
- `expiration_ms` (optional): Expiration time in milliseconds

### GET /urls/{short_code}

Redirect to the original URL.

**Response:**
- Status: 301 (Moved Permanently)
- Location header: Original URL

## Example Usage

### Create a short URL:
```bash
curl -X POST http://localhost:8080/urls \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.google.com",
    "alias": "google",
    "expiration_ms": 3600000
  }'
```

### Access the short URL:
```bash
curl -I http://localhost:8080/urls/google
```

## Notes

- Currently uses in-memory storage (data is lost on server restart)
- Short codes are 5 characters long when auto-generated
- Expired URLs are automatically cleaned up when accessed
- Custom aliases must be unique

## Error Handling

The API uses a centralized error handling system where each error type carries its own HTTP status code. This ensures consistent error responses across all endpoints.

### Error Types and Status Codes

| Error Type | HTTP Status | Description |
|------------|-------------|-------------|
| Invalid URL format | **400** | Bad Request |
| Invalid alias format | **400** | Bad Request |
| Alias already exists | **409** | Conflict |
| Short code not found | **404** | Not Found |
| Short code expired | **404** | Not Found |
| Server errors | **500** | Internal Server Error |

### Benefits of This Approach

- **Centralized**: All error definitions are in one place
- **Consistent**: Same error always returns same status code
- **Maintainable**: Easy to add new error types
- **Clean Handlers**: Handlers don't need to know about status codes

## Testing

### Running Tests

The project includes comprehensive unit tests and integration tests for all components.

#### Quick Test Run
```bash
chmod +x run_tests.sh
./run_tests.sh
```

#### Individual Test Categories
```bash
# Unit tests for services
go test -v ./tests/unit/services/...

# Unit tests for handlers
go test -v ./tests/unit/handlers/...

# Unit tests for models
go test -v ./tests/unit/models/...

# Integration tests
go test -v ./tests/integration/...

# All tests with coverage
go test -v -coverprofile=coverage.out ./...
```

### Test Structure

- **Unit Tests**: Test individual components in isolation
  - `services/`: Business logic tests
  - `handlers/`: HTTP handler tests
  - `models/`: Data structure and error tests

- **Integration Tests**: Test complete API workflows
  - Happy path scenarios
  - Error handling scenarios
  - Edge cases

### Test Coverage

The test suite covers:
- ✅ URL validation (format, scheme, normalization)
- ✅ Alias validation (length, characters, uniqueness)
- ✅ Short code generation and storage
- ✅ URL expiration handling
- ✅ Error responses and HTTP status codes
- ✅ Complete API workflows

### Manual API Testing

You can also test the API manually using the provided test script:

```bash
chmod +x test_errors.sh
./test_errors.sh
```
