[![CI/CD Pipeline](https://github.com/coelhocaique/url-shortener-api/actions/workflows/ci-cd.yml/badge.svg?branch=main)](https://github.com/coelhocaique/url-shortener-api/actions/workflows/ci-cd.yml)

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

### Option 1: Using Docker (Recommended)

1. Make sure you have Docker and Docker Compose installed
2. Navigate to the project directory:
   ```bash
   cd url-shortener-api
   ```

3. Build and run with Docker Compose:
   ```bash
   docker-compose up --build
   ```

4. The API will be available at `http://localhost:8080`

### Option 2: Local Development

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

### Docker Commands

```bash
# Build the Docker image
docker build -t url-shortener-api .

# Run the container
docker run -p 8080:8080 url-shortener-api

# Run with custom port
docker run -p 3000:8080 -e PORT=8080 url-shortener-api

# Run in detached mode
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

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

## CI/CD Pipeline

This project includes a comprehensive GitHub Actions CI/CD pipeline that runs on every push to the main branch and pull requests.

### Pipeline Features

- **Automated Testing**: Runs unit and integration tests with coverage reporting
- **Docker Build**: Builds and pushes Docker images to Docker Hub
- **Security Scanning**: Performs vulnerability scanning with Trivy
- **Multi-stage Pipeline**: Separate jobs for testing, building, security scanning, and deployment

### Pipeline Jobs

1. **Test Job**
   - Runs Go tests with race detection
   - Generates coverage reports
   - Uploads coverage artifacts

2. **Build Job**
   - Builds Docker image using multi-stage build
   - Pushes to Docker Hub (on main branch)
   - Uses Docker layer caching for faster builds

3. **Security Scan Job**
   - Scans Docker image for vulnerabilities
   - Uploads results to GitHub Security tab
   - Only runs on main branch pushes

4. **Deploy Job**
   - Deploys to production environment
   - Protected by environment rules
   - Only runs on main branch pushes

### Required Secrets

To enable the CI/CD pipeline, add these secrets to your GitHub repository:

- `DOCKER_USERNAME`: Your Docker Hub username
- `DOCKER_PASSWORD`: Your Docker Hub password or access token

### Setting up Secrets

1. Go to your GitHub repository
2. Navigate to Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Add the required secrets

### Pipeline Triggers

- **Push to main**: Runs full pipeline including deployment
- **Pull Request**: Runs testing and building (no deployment)
- **Manual trigger**: Can be triggered manually from Actions tab

### Monitoring

- View pipeline status in the Actions tab
- Coverage reports are uploaded as artifacts
- Security scan results appear in the Security tab
- Build logs are available for debugging
