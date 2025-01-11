# Metadata Management Service

A RESTful API service built with Go, Fiber, and PostgreSQL for managing metadata with JSON values. The service implements CRUD operations with support for partial updates and follows clean architecture principles.

## Features

- Complete CRUD operations for metadata entries
- JSON value storage using PostgreSQL JSONB
- Key validation (alphanumeric, underscore, hyphen, and dot allowed)
- Transaction support for all database operations
- CORS enabled
- Clean architecture with Factory, Service, and Repository patterns
- Comprehensive error handling and HTTP responses

## Architecture

The application follows a clean, layered architecture:

1. **API Layer** (Handler)
   - HTTP request/response handling
   - Input validation
   - Route management

2. **Business Layer** (Service)
   - Business logic
   - Key format validation
   - Operation orchestration

3. **Data Access Layer** (Repository)
   - Database operations
   - Transaction management
   - Data persistence

4. **Factory Layer**
   - Dependency injection
   - Component creation
   - Instance management

## Prerequisites

- Go 1.19 or later
- PostgreSQL 15 or later
- Docker and Docker Compose (for local development)

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd metadata-service
```

2. Start PostgreSQL using Docker:
```bash
docker-compose up -d
```

3. Install dependencies:
```bash
go mod tidy
```

4. Run the application:
```bash
go run .
```

The server will start on `http://localhost:3000`.

## API Endpoints

### Create Entry
```bash
POST /api/metadata
{
    "my_key": "config1",
    "my_value": {
        "environment": "production",
        "max_connections": 100
    }
}
```

### Get Entry
```bash
GET /api/metadata/:key
```

### Update Entry
```bash
PUT /api/metadata/:key
{
    "my_value": {
        "environment": "staging",
        "max_connections": 50
    }
}
```

### Patch Update Entry
```bash
PATCH /api/metadata/:key
{
    "environment": "development"
}
```

### Delete Entry
```bash
DELETE /api/metadata/:key
```

### Delete All Entries
```bash
DELETE /api/metadata
```

### Upsert Entry
```bash
PUT /api/metadata
{
    "my_key": "config1",
    "my_value": {
        "environment": "production",
        "max_connections": 100
    }
}
```

## Response Formats

### Success Response
```json
{
    "message": "Operation successful",
    "data": {
        "my_key": "config1",
        "my_value": {
            "environment": "production",
            "max_connections": 100
        }
    }
}
```

### Error Response
```json
{
    "error": "Error message",
    "details": "Detailed error information"
}
```

## Project Structure

```
├── main.go           # Application entry point and route setup
├── handler.go        # HTTP request handlers
├── service.go        # Business logic layer
├── repository.go     # Data access layer
├── factory.go        # Dependency injection
├── db.go            # Database models and connection
├── docker-compose.yml    # Docker configuration
└── init.sql         # Database initialization
```

## Key Validation Rules

Keys must follow these rules:
- Can contain alphanumeric characters (a-z, A-Z, 0-9)
- Can contain underscore (_)
- Can contain hyphen (-)
- Can contain dot (.)
- Cannot be empty

## Development

### Adding New Features

1. Define the interface in the appropriate layer
2. Implement the interface
3. Update the factory if needed
4. Add necessary tests
5. Update documentation

### Running Tests
```bash
go test ./...
```

## Database Schema

```sql
CREATE TABLE metadata_table (
    my_key TEXT PRIMARY KEY,
    my_value JSONB NOT NULL
);

CREATE INDEX idx_metadata_key ON metadata_table(my_key);
CREATE INDEX idx_metadata_value ON metadata_table USING gin(my_value);
```

## Error Handling

The service provides detailed error responses for various scenarios:
- Invalid key format
- Duplicate keys
- Missing required fields
- Invalid JSON format
- Not found errors
- Server errors