# CMS Service Structure

## Overview

The CMS service provides a flexible, production-ready client for Strapi 5 with Redis caching, following best practices for Go services.

## Directory Structure

```
services/cms/
│
├── 📄 client.go              - Main CMS client (HTTP, auth, caching)
├── 📄 types.go               - Query options and error types
├── 📄 cache.go               - Redis caching layer with interface
│
├── 📁 models/                - Content type models (one per file)
│   ├── 📄 strapi.go          - Core Strapi 5 response structures
│   ├── 📄 fields.go          - Common field types (media, relations)
│   ├── 📄 article.go         - Example: Article collection type
│   ├── 📄 homepage.go        - Example: Homepage single type
│   └── 📄 README.md          - Model creation guide
│
├── 📄 example_service.go     - Example services (ArticleService, HomepageService)
│
├── 📄 README.md              - Complete usage documentation
└── 📄 STRUCTURE.md           - This file
```

## Component Responsibilities

### client.go
- HTTP client with proper timeouts and connection pooling
- Authentication (Bearer token)
- Context support for cancellation
- Cache integration
- URL building for Strapi 5 API
- Error handling with Strapi error responses

### types.go
- `CollectionOptions` - Query parameters for collections
- `ItemOptions` - Query parameters for single items
- `CacheOptions` - Cache configuration per request
- `RequestError` - Structured error type

### cache.go
- `Cache` interface - Abstraction for caching
- `RedisCache` - Redis implementation
- `NoOpCache` - Fallback when Redis unavailable
- `CacheKeyBuilder` - Consistent cache key generation

### models/
Each model file contains:
- Struct definition matching Strapi content type
- Type aliases for response types
- Uses common field types from `fields.go`

### Service Files
Each service file (e.g., `articles_service.go`) contains:
- Service struct with client reference
- Methods: `GetAll()`, `GetByID()`, `GetBySlug()`, etc.
- Cache management per resource
- Response unmarshaling

## Data Flow

```
┌─────────────────────────────────────────────────────────────┐
│                        Fiber Handler                        │
│                     (in main.go)                            │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          │ calls method
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    Resource Service                         │
│              (e.g., ArticleService)                         │
│  • GetAll(opts, useCache)                                   │
│  • GetByID(id, opts, useCache)                              │
│  • GetBySlug(slug, populate, useCache)                      │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          │ uses
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                      CMS Client                             │
│  • GetCollection(endpoint, opts, cacheOpts)                 │
│  • GetItem(endpoint, id, opts, cacheOpts)                   │
│  • GetSingle(endpoint, opts, cacheOpts)                     │
└──────────┬──────────────────────────────────┬───────────────┘
           │                                  │
           │ checks cache                     │ makes request
           ▼                                  ▼
┌─────────────────────┐          ┌──────────────────────────┐
│   Redis Cache       │          │   Strapi 5 API           │
│  • Get()            │          │  (HTTP Request)          │
│  • Set()            │          │  • Authorization header  │
│  • Delete()         │          │  • Query parameters      │
└─────────────────────┘          └──────────────────────────┘
           │                                  │
           │                                  │ returns JSON
           └──────────────┬───────────────────┘
                          │
                          │ unmarshals to
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                     Model Struct                            │
│        StrapiResponse[YourContentType]                      │
│                                                             │
│  {                                                          │
│    data: {                                                  │
│      id: 1,                                                 │
│      attributes: { ...your fields... }                     │
│    }                                                        │
│  }                                                          │
└─────────────────────────────────────────────────────────────┘
```

## Creating a New Resource

### Step 1: Create Model
```bash
# Create models/product.go
```

```go
package models

type Product struct {
    Name  string `json:"name"`
    Price float64 `json:"price"`
    // ... more fields
}

type ProductResponse = StrapiResponse[Product]
type ProductCollectionResponse = StrapiCollectionResponse[Product]
```

### Step 2: Create Service
```bash
# Create products_service.go
```

```go
package cms

import (
    "api-gateway/services/cms/models"
    "context"
)

type ProductService struct {
    client   Client
    endpoint string
    cacheTTL time.Duration
}

func NewProductService(client Client) *ProductService {
    return &ProductService{
        client:   client,
        endpoint: "products",
        cacheTTL: 10 * time.Minute,
    }
}

func (s *ProductService) GetAll(ctx context.Context, opts CollectionOptions, useCache bool) (*models.ProductCollectionResponse, error) {
    // Implementation similar to ArticleService
}

func (s *ProductService) GetByID(ctx context.Context, id string, opts ItemOptions, useCache bool) (*models.ProductResponse, error) {
    // Implementation similar to ArticleService
}
```

### Step 3: Initialize in main.go
```go
productService := cms.NewProductService(cmsClient)
```

### Step 4: Add Handler
```go
app.Get("/api/cms/products", func(c *fiber.Ctx) error {
    ctx := c.Context()

    opts := cms.CollectionOptions{
        Page:     1,
        PageSize: 10,
        Populate: "*",
    }

    products, err := productService.GetAll(ctx, opts, true)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(products)
})
```

## Key Features

### ✅ Type Safety
- Generic types for Strapi responses
- Compile-time type checking
- IDE autocomplete support

### ✅ Caching
- Redis-backed caching
- Per-request cache control
- Automatic key generation
- Pattern-based invalidation
- Graceful fallback

### ✅ Best Practices
- Context for cancellation
- Proper HTTP client configuration
- Connection pooling
- Timeout handling
- Structured error handling
- Interface-based design

### ✅ Flexibility
- Support for all Strapi 5 query options
- Pagination, filtering, sorting, population
- Localization support
- Custom cache keys
- Service-per-resource pattern

## Environment Variables

```env
# Strapi Configuration
CMS_SERVICE_URL=http://localhost:1337/api
CMS_SERVICE_TOKEN=your-api-token

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
```

## Testing Without Redis

The client gracefully handles Redis unavailability:

```go
// If Redis connection fails, NoOpCache is used automatically
// No code changes needed - caching simply disabled
```

## Next Steps

1. ✏️ Create model files for your Strapi content types in `models/`
2. ✏️ Create service files for each resource
3. ✏️ Initialize services in `main.go`
4. ✏️ Add Fiber handlers for your endpoints
5. 🚀 Deploy with Redis for production caching

## Resources

- [models/README.md](models/README.md) - Model creation guide
- [README.md](README.md) - Complete usage documentation
- [example_service.go](example_service.go) - Service implementation examples
