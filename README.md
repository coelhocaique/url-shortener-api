[![CI/CD Pipeline](https://github.com/coelhocaique/url-shortener-api/actions/workflows/ci-cd.yml/badge.svg?branch=main)](https://github.com/coelhocaique/url-shortener-api/actions/workflows/ci-cd.yml)

# URL Shortener API

A modular URL shortener API built with Go and Gin framework, following clean architecture principles. 

The project implements the design below

![alt text](https://github.com/coelhocaique/url-shortener-api/blob/main/design.png?raw=true)


## Features

- Create short URLs with optional custom aliases
- Set expiration times for URLs
- Automatic redirect to original URLs
- MongoDB persistence with automatic TTL cleanup
- User-specific URL management
- Modular architecture with separation of concerns
- Comprehensive indexing for optimal performance

## Setup

### Option 1: Using Docker (Recommended)

1. Make sure you have Docker and Docker Compose installed
2. Navigate to the project directory:
   ```bash
   cd url-shortener-api
   ```

3. Build and run with Docker Compose (includes MongoDB):
   ```bash
   docker-compose up --build
   ```

4. The API will be available at `http://localhost:8080`
5. MongoDB will be available at `localhost:27017`

### Option 2: Local Development

1. Make sure you have Go 1.21+ installed
2. Install and start MongoDB locally:
   ```bash
   # On macOS with Homebrew
   brew install mongodb-community
   brew services start mongodb-community
   
   # On Ubuntu/Debian
   sudo apt-get install mongodb
   sudo systemctl start mongod
   
   # Or use Docker for MongoDB only
   docker run -d -p 27017:27017 --name mongodb mongo:7.0
   ```

3. Navigate to the project directory:
   ```bash
   cd url-shortener-api
   ```

4. Install dependencies:
   ```bash
   go mod tidy
   ```

5. Set environment variables (optional):
   ```bash
   export MONGO_URI=mongodb://localhost:27017
   export DATABASE_NAME=url_shortener
   export PORT=8080
   ```

6. Run the server:
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
   "expiration_ms": "1000000",
   "user_id": "user123"
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
- `user_id` (optional): User identifier for URL ownership

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
    "expiration_ms": 3600000,
    "user_id": "user123"
  }'
```

### Access the short URL:
```bash
curl -I http://localhost:8080/urls/google
```

## Notes

- Uses MongoDB for persistent storage with automatic TTL cleanup
- Short codes are 5 characters long when auto-generated
- Expired URLs are automatically cleaned up by MongoDB TTL indexes
- Custom aliases must be unique across all users
- User-specific URL management with user_id field
- Comprehensive indexing for optimal query performance

## Data Model

The application uses MongoDB to store URL mappings with the following document structure:

```json
{
  "_id": "ObjectId",
  "original_url": "https://www.google.com",
  "short_url": "abc123",
  "expiration_timestamp": "2024-12-31T23:59:59Z",
  "alias": "google",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "user_id": "user123"
}
```

### Indexes

The following indexes are automatically created for optimal performance:

- **short_url**: Unique index for fast lookups
- **alias**: Unique sparse index for custom aliases
- **user_id**: Index for user-specific queries
- **expiration_timestamp**: TTL index for automatic cleanup

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

The project includes comprehensive unit tests and integration tests for all components. Tests now use MongoDB for persistence.

#### Quick Test Run (with automatic MongoDB setup)
```bash
chmod +x run_tests_mongodb.sh
./run_tests_mongodb.sh
```

#### Manual Test Run (requires MongoDB to be running)
```bash
chmod +x run_tests.sh
./run_tests.sh
```

#### Individual Test Categories
```bash
# Make sure MongoDB is running first
docker run -d -p 27017:27017 --name test-mongodb mongo:7.0

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

# Cleanup
docker stop test-mongodb && docker rm test-mongodb
```

### Test Structure

- **Unit Tests**: Test individual components in isolation
  - `services/`: Business logic tests with MongoDB integration
  - `handlers/`: HTTP handler tests with mocked services
  - `models/`: Data structure and error tests

- **Integration Tests**: Test complete API workflows with MongoDB
  - Happy path scenarios
  - Error handling scenarios
  - Edge cases
  - Database persistence and retrieval

- **Test Utilities**: MongoDB test setup and utilities
  - `testutils/`: MongoDB connection and cleanup utilities

### Test Coverage

The test suite covers:
- ✅ URL validation (format, scheme, normalization)
- ✅ Alias validation (length, characters, uniqueness)
- ✅ Short code generation and MongoDB storage
- ✅ URL expiration handling with TTL indexes
- ✅ Error responses and HTTP status codes
- ✅ Complete API workflows with database persistence
- ✅ MongoDB-specific operations (GetByAlias, GetByUserID)
- ✅ Database indexing and performance

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
