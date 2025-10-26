package models

import "time"

// Strapi 5 API Response structures

// StrapiResponse represents a single item response from Strapi 5
type StrapiResponse[T any] struct {
	Data *StrapiData[T] `json:"data"`
	Meta *StrapiMeta    `json:"meta,omitempty"`
}

// StrapiCollectionResponse represents a collection response from Strapi 5
type StrapiCollectionResponse[T any] struct {
	Data []StrapiData[T] `json:"data"`
	Meta *StrapiMeta     `json:"meta,omitempty"`
}

// StrapiData wraps the actual content with metadata
type StrapiData[T any] struct {
	ID          int        `json:"id"`
	Attributes  T          `json:"attributes"`
	CreatedAt   time.Time  `json:"createdAt,omitempty"`
	UpdatedAt   time.Time  `json:"updatedAt,omitempty"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
}

// StrapiMeta contains pagination and other metadata
type StrapiMeta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination contains pagination information
type Pagination struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	PageCount int `json:"pageCount"`
	Total     int `json:"total"`
}

// StrapiError represents an error response from Strapi
type StrapiError struct {
	Error struct {
		Status  int                    `json:"status"`
		Name    string                 `json:"name"`
		Message string                 `json:"message"`
		Details map[string]interface{} `json:"details,omitempty"`
	} `json:"error"`
}
