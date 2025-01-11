# Understanding Clean Architecture Layers

This document explains the different architectural layers used in our metadata management service. The application follows clean architecture principles, separating concerns into distinct layers.

## Overview

```
Client Request
     ↓
[API Layer (Handler)]
     ↓
[Business Layer (Service)]
     ↓
[Data Access Layer (Repository)]
     ↓
[Database]
```

## Factory Pattern

The Factory pattern is responsible for object creation and dependency management.

### Purpose
- Centralizes object creation logic
- Manages dependencies
- Makes testing easier through dependency injection
- Allows easy swapping of implementations

### Example in Our Code
```go
type MetadataFactory struct {
    db *gorm.DB
}

func (f *MetadataFactory) CreateService() *MetadataService {
    return NewMetadataService(f.CreateRepository())
}
```

### Usage
```go
factory := NewMetadataFactory(db)
service := factory.CreateService()
```

## Service Layer (Business Layer)

The Service layer contains business logic and rules for the application.

### Purpose
- Implements business rules and validations
- Orchestrates operations between layers
- Keeps business logic separate from technical implementation
- Maintains application's core functionality

### Example in Our Code
```go
type MetadataService struct {
    repo MetadataRepository
}

func (s *MetadataService) validateKey(key string) error {
    if !keyPattern.MatchString(key) {
        return fmt.Errorf("invalid key format...")
    }
    return nil
}
```

### Responsibilities
1. Key format validation
2. Business rule enforcement
3. Operation orchestration
4. Error handling for business rules

## Repository Pattern (Data Access Layer)

The Repository pattern abstracts data persistence operations.

### Purpose
- Abstracts database operations
- Provides collection-like interface for data access
- Manages transactions
- Encapsulates data access logic

### Example in Our Code
```go
type MetadataRepository interface {
    Create(entry *MetadataEntry) error
    Get(key string) (*MetadataEntry, error)
    Update(entry *MetadataEntry) error
    Delete(key string) error
}

func (r *metadataRepository) Create(entry *MetadataEntry) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        return tx.Create(entry).Error
    })
}
```

### Responsibilities
1. CRUD operations
2. Transaction management
3. Data persistence
4. Query execution

## Data Access Layer

The Data Access Layer handles the actual data storage and retrieval.

### Purpose
- Defines database models
- Manages database connections
- Handles data mapping
- Provides database abstraction

### Example in Our Code
```go
type MetadataEntry struct {
    MyKey   string          `gorm:"column:my_key;primaryKey"`
    MyValue json.RawMessage `gorm:"column:my_value;type:jsonb"`
}

func InitDB() (*gorm.DB, error) {
    dsn := "host=localhost user=postgres..."
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
```

### Responsibilities
1. Database schema definition
2. Connection management
3. Data type mapping
4. Database configuration

## Flow of a Request

Let's follow a request through all layers:

1. **Client Makes Request**
   ```http
   POST /api/metadata
   {"my_key": "config1", "my_value": {"setting": "value"}}
   ```

2. **Handler Layer (API)**
   ```go
   func (h *MetadataHandler) Create(c *fiber.Ctx) error {
       // Parses request
       // Validates format
       // Calls service layer
   }
   ```

3. **Service Layer (Business)**
   ```go
   func (s *MetadataService) Create(entry *MetadataEntry) error {
       // Validates business rules
       // Calls repository layer
   }
   ```

4. **Repository Layer (Data Access)**
   ```go
   func (r *metadataRepository) Create(entry *MetadataEntry) error {
       // Manages transaction
       // Performs database operation
   }
   ```

## Benefits of This Architecture

1. **Separation of Concerns**
   - Each layer has a specific responsibility
   - Changes in one layer don't affect others
   - Clear boundaries between components

2. **Maintainability**
   - Easy to modify individual components
   - Clear structure for new developers
   - Isolated testing possible

3. **Testability**
   - Each layer can be tested independently
   - Easy to mock dependencies
   - Clear interfaces for testing

4. **Scalability**
   - Easy to modify implementations
   - Can change database without affecting business logic
   - Can add new features without disturbing existing ones

## Common Operations Example

### Creating a New Entry
```
HTTP Request
     ↓
Handler validates JSON format
     ↓
Service validates business rules
     ↓
Repository manages transaction
     ↓
Data access persists to database
```

### Best Practices

1. **Factory Pattern**
   - Use for complex object creation
   - Keep creation logic centralized
   - Enable dependency injection

2. **Service Layer**
   - Keep business logic isolated
   - Validate business rules
   - Handle business exceptions

3. **Repository Pattern**
   - Use transactions for data consistency
   - Abstract database operations
   - Keep data access logic contained

4. **Data Access Layer**
   - Define clear data models
   - Handle database connections
   - Manage data mappings

## Testing Considerations

Each layer can be tested independently:

```go
// Service Layer Test
func TestMetadataService_ValidateKey(t *testing.T) {
    service := NewMetadataService(mockRepository)
    err := service.validateKey("invalid@key")
    assert.Error(t, err)
}

// Repository Layer Test
func TestMetadataRepository_Create(t *testing.T) {
    repo := NewMetadataRepository(mockDB)
    err := repo.Create(testEntry)
    assert.NoError(t, err)
}
```