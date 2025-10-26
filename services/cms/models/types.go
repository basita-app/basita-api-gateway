package models

import (
	"time"
)

// CollectionOptions contains query options for fetching collections
type CollectionOptions struct {
	Page     int               // Page number for pagination
	PageSize int               // Number of items per page
	Populate string            // Fields to populate (e.g., "author,categories")
	Filters  map[string]string // Filter criteria
	Sort     []string          // Sort fields (e.g., ["title:asc", "publishedAt:desc"])
	Locale   string            // Locale for internationalization
	Fields   []string          // Specific fields to return
}

// ItemOptions contains query options for fetching a single item
type ItemOptions struct {
	Populate string   // Fields to populate
	Locale   string   // Locale for internationalization
	Fields   []string // Specific fields to return
}

// CacheOptions controls caching behavior
type CacheOptions struct {
	Enabled bool          // Enable/disable cache for this request
	TTL     time.Duration // Time to live for cached data
	Key     string        // Custom cache key (optional, auto-generated if empty)
}

// RequestError represents an error that occurred during a request
type RequestError struct {
	StatusCode int
	Message    string
	StrapiErr  *StrapiError
}

func (e *RequestError) Error() string {
	if e.StrapiErr != nil {
		return e.StrapiErr.Error.Message
	}
	return e.Message
}
