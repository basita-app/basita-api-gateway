# CMS Service Client

A flexible, production-ready Go client for interacting with Strapi 5 CMS with built-in Redis caching support.

## Features

- ✅ Full Strapi 5 API support (Collections, Single Items, Single Types)
- ✅ Redis caching with configurable TTL
- ✅ Type-safe responses with generics
- ✅ Context support for cancellation and timeouts
- ✅ Proper error handling with Strapi error responses
- ✅ Authentication with Bearer tokens
- ✅ Configurable HTTP client with connection pooling
- ✅ Cache invalidation support
- ✅ Flexible query options (pagination, filters, population, sorting)

## Quick Start

### 1. Create a Model

Create a file in `models/` for your Strapi content type:

```go
// models/product.go
package models

type Product struct {
    Name        string  `json:"name"`
    Slug        string  `json:"slug"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}

type ProductResponse = StrapiResponse[Product]
type ProductCollectionResponse = StrapiCollectionResponse[Product]
```

### 2. Create a Service

Create a service file for your resource:

```go
// products_service.go
package cms

import (
    "api-gateway/services/cms/models"
    "context"
    "encoding/json"
    "time"
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
    var cacheOpts *CacheOptions
    if useCache {
        cacheOpts = &CacheOptions{Enabled: true, TTL: s.cacheTTL}
    }

    data, err := s.client.GetCollection(ctx, s.endpoint, opts, cacheOpts)
    if err != nil {
        return nil, err
    }

    var response models.ProductCollectionResponse
    if err := json.Unmarshal(data, &response); err != nil {
        return nil, &RequestError{Message: "failed to unmarshal products: " + err.Error()}
    }

    return &response, nil
}

func (s *ProductService) GetByID(ctx context.Context, id string, opts ItemOptions, useCache bool) (*models.ProductResponse, error) {
    var cacheOpts *CacheOptions
    if useCache {
        cacheOpts = &CacheOptions{Enabled: true, TTL: s.cacheTTL}
    }

    data, err := s.client.GetItem(ctx, s.endpoint, id, opts, cacheOpts)
    if err != nil {
        return nil, err
    }

    var response models.ProductResponse
    if err := json.Unmarshal(data, &response); err != nil {
        return nil, &RequestError{Message: "failed to unmarshal product: " + err.Error()}
    }

    return &response, nil
}
```

### 3. Initialize in main.go

```go
// The CMS client is already initialized in main.go with Redis caching
// Just create your service:
productService := cms.NewProductService(cmsClient)
```

### 4. Create Handlers

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

## Query Options

### Collection Options

```go
opts := cms.CollectionOptions{
    Page:     1,
    PageSize: 10,
    Populate: "*",
    Locale:   "en",
    Sort:     []string{"createdAt:desc"},
    Fields:   []string{"title", "slug"},
    Filters: map[string]string{
        "filters[status][$eq]": "published",
    },
}
```

### Item Options

```go
opts := cms.ItemOptions{
    Populate: "*",
    Locale:   "en",
    Fields:   []string{"title", "content"},
}
```

## Caching

```go
// Enable caching
cacheOpts := &cms.CacheOptions{
    Enabled: true,
    TTL:     15 * time.Minute,
}

// Disable caching (pass nil)
data, err := client.GetCollection(ctx, endpoint, opts, nil)

// Invalidate cache
err := service.InvalidateCache(ctx) // Invalidates all cache for this service
```

## Environment Variables

```env
# CMS Configuration
CMS_SERVICE_URL=http://localhost:1337/api
CMS_SERVICE_TOKEN=your-strapi-api-token

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
```

## Documentation

- **[docs/cms/STRUCTURE.md](../../docs/cms/STRUCTURE.md)** - Detailed architecture and data flow
- **[docs/cms/MODELS.md](../../docs/cms/MODELS.md)** - Complete guide for creating models
- **[pkg/cache/README.md](../../pkg/cache/README.md)** - Cache package documentation

## File Structure

```
services/cms/
├── client.go              # Main CMS client implementation
├── types.go               # Query options and error types
├── cache_key_builder.go   # CMS-specific cache key generation
├── models/                # Content type models (create one per Strapi type)
│   ├── strapi.go          # Core Strapi 5 response structures
│   └── fields.go          # Common field types (media, relations, components)
└── README.md              # This file
```

## Common Patterns

### Single Type (e.g., Homepage, Settings)

```go
type Homepage struct {
    Title string `json:"title"`
}

func (s *HomepageService) Get(ctx context.Context, populate string, useCache bool) (*models.HomepageResponse, error) {
    opts := ItemOptions{Populate: populate}
    cacheOpts := &CacheOptions{Enabled: useCache, TTL: 30 * time.Minute}

    data, err := s.client.GetSingle(ctx, s.endpoint, opts, cacheOpts)
    if err != nil {
        return nil, err
    }

    var response models.HomepageResponse
    json.Unmarshal(data, &response)
    return &response, nil
}
```

### Get By Slug

```go
func (s *ArticleService) GetBySlug(ctx context.Context, slug string, useCache bool) (*models.ArticleResponse, error) {
    opts := CollectionOptions{
        Filters: map[string]string{"filters[slug][$eq]": slug},
        PageSize: 1,
    }

    data, err := s.client.GetCollection(ctx, s.endpoint, opts, cacheOpts)
    // ... unmarshal and return first item
}
```

## Best Practices

1. **One model per file** in `models/` directory
2. **One service per resource** following the naming pattern `{resource}_service.go`
3. **Always use context** for proper cancellation
4. **Enable caching** for read-heavy endpoints
5. **Use type aliases** for convenience (e.g., `ProductResponse`)
6. **Handle errors properly** with RequestError type
7. **Set appropriate cache TTLs** based on data volatility

## Dependencies

```bash
go get github.com/redis/go-redis/v9
go get github.com/gofiber/fiber/v2
```
