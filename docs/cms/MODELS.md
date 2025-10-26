# CMS Models

This directory contains all Strapi content type models. Each content type should have its own file.

## File Structure

```
models/
├── strapi.go       # Core Strapi 5 response structures (don't modify)
├── fields.go       # Common Strapi field types (media, relations, components)
├── article.go      # Example: Article collection type
├── homepage.go     # Example: Homepage single type
└── your_model.go   # Your custom models
```

## Creating a New Model

### For Collection Types

Create a new file `models/your_resource.go`:

```go
package models

import "time"

// YourResource represents the attributes for your Strapi collection type
type YourResource struct {
    Title       string     `json:"title"`
    Slug        string     `json:"slug"`
    Description string     `json:"description,omitempty"`
    PublishedAt *time.Time `json:"publishedAt,omitempty"`

    // Relations
    Author     *RelationField[Author]          `json:"author,omitempty"`
    Categories *RelationCollectionField[Category] `json:"categories,omitempty"`

    // Media
    Image *MediaField `json:"image,omitempty"`

    // Add more fields matching your Strapi schema
}

// Type aliases for convenience
type YourResourceResponse = StrapiResponse[YourResource]
type YourResourceCollectionResponse = StrapiCollectionResponse[YourResource]
```

### For Single Types

Create a new file `models/your_single_type.go`:

```go
package models

// YourSingleType represents the attributes for a Strapi single type
type YourSingleType struct {
    Title       string `json:"title"`
    Content     string `json:"content"`

    // Media
    HeroImage *MediaField `json:"heroImage,omitempty"`

    // Components
    Sections *ComponentCollectionField[Section] `json:"sections,omitempty"`

    // Add more fields matching your Strapi schema
}

// Type alias for convenience
type YourSingleTypeResponse = StrapiResponse[YourSingleType]
```

## Using Common Field Types

### Media Fields

```go
// Single media field
Image *MediaField `json:"image,omitempty"`

// Access URL: response.Data.Attributes.Image.Data.Attributes.URL
```

### Relation Fields

```go
// Single relation
Author *RelationField[Author] `json:"author,omitempty"`

// Multiple relations
Categories *RelationCollectionField[Category] `json:"categories,omitempty"`

// Access: response.Data.Attributes.Author.Data.Attributes.Name
```

### Component Fields

```go
// Single component
Hero *ComponentField[HeroSection] `json:"hero,omitempty"`

// Repeatable component
Features *ComponentCollectionField[Feature] `json:"features,omitempty"`
```

## Nested Models Example

If you have relations, create models for those too:

```go
// models/author.go
package models

type Author struct {
    Name   string      `json:"name"`
    Bio    string      `json:"bio,omitempty"`
    Avatar *MediaField `json:"avatar,omitempty"`
}

type AuthorResponse = StrapiResponse[Author]
type AuthorCollectionResponse = StrapiCollectionResponse[Author]
```

```go
// models/category.go
package models

type Category struct {
    Name string `json:"name"`
    Slug string `json:"slug"`
}

type CategoryResponse = StrapiResponse[Category]
type CategoryCollectionResponse = StrapiCollectionResponse[Category]
```

## Components Example

For Strapi components, create them in the same models directory:

```go
// models/hero_section.go
package models

type HeroSection struct {
    Title    string      `json:"title"`
    Subtitle string      `json:"subtitle,omitempty"`
    Image    *MediaField `json:"image,omitempty"`
    CTAText  string      `json:"ctaText,omitempty"`
    CTALink  string      `json:"ctaLink,omitempty"`
}
```

## Best Practices

1. **One model per file**: Each content type gets its own file
2. **Match Strapi schema exactly**: Field names and types should match your Strapi content type
3. **Use json tags**: Always include `json:"fieldName"` tags
4. **Use omitempty**: For optional fields, use `omitempty` in json tags
5. **Create type aliases**: Add convenience type aliases for responses
6. **Document your fields**: Add comments for complex or non-obvious fields
7. **Reuse common types**: Use `MediaField`, `RelationField`, etc. from `fields.go`
8. **Keep models simple**: Models should only define structure, no business logic

## Strapi 5 Field Mapping

| Strapi Type | Go Type |
|-------------|---------|
| Text (short/long) | `string` |
| Rich Text | `string` |
| Number (integer) | `int` or `int64` |
| Number (float) | `float64` |
| Boolean | `bool` |
| Date | `time.Time` |
| DateTime | `time.Time` |
| Email | `string` |
| Enumeration | `string` |
| Media (single) | `*MediaField` |
| Media (multiple) | `*MediaCollectionField` |
| Relation (single) | `*RelationField[T]` |
| Relation (multiple) | `*RelationCollectionField[T]` |
| Component (single) | `*ComponentField[T]` |
| Component (repeatable) | `*ComponentCollectionField[T]` |
| Dynamic Zone | `[]DynamicZoneItem` |
| JSON | `map[string]interface{}` |

## Example: Complete Product Model

```go
// models/product.go
package models

import "time"

type Product struct {
    // Basic fields
    Name        string  `json:"name"`
    Slug        string  `json:"slug"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    SKU         string  `json:"sku"`
    InStock     bool    `json:"inStock"`

    // Media
    MainImage *MediaField            `json:"mainImage,omitempty"`
    Gallery   *MediaCollectionField  `json:"gallery,omitempty"`

    // Relations
    Category     *RelationField[Category]              `json:"category,omitempty"`
    Tags         *RelationCollectionField[Tag]         `json:"tags,omitempty"`
    RelatedItems *RelationCollectionField[Product]     `json:"relatedItems,omitempty"`

    // Components
    Specifications *ComponentCollectionField[Specification] `json:"specifications,omitempty"`

    // Timestamps
    PublishedAt *time.Time `json:"publishedAt,omitempty"`
}

type ProductResponse = StrapiResponse[Product]
type ProductCollectionResponse = StrapiCollectionResponse[Product]
```

## Getting Help

- Check [Strapi 5 REST API docs](https://docs.strapi.io/dev-docs/api/rest) for API structure
- Review `example_service.go` for usage examples
- See `fields.go` for available field type helpers
