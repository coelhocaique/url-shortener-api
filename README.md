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
├── models/          # Data structures and interfaces
├── routes/          # Route definitions
├── services/        # Business logic
│   ├── factory.go           # Service dependency injection
│   ├── short_code_generator.go  # Short code generation
│   ├── url_service.go       # Main URL service
│   ├── url_storage.go       # URL storage operations
│   └── url_validator.go     # URL validation
├── main.go         # Application entry point
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
