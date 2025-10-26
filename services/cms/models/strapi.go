package models

import (
	"encoding/json"
	"fmt"
	"time"
)

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
// In Strapi 5, attributes are embedded directly in the data object (not nested)
type StrapiData[T any] struct {
	ID          int        `json:"id"`
	DocumentID  string     `json:"documentId,omitempty"`
	Attributes  T          `json:"-"` // We'll handle this manually
	CreatedAt   time.Time  `json:"createdAt,omitempty"`
	UpdatedAt   time.Time  `json:"updatedAt,omitempty"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	Locale      string     `json:"locale,omitempty"`
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

// UnmarshalJSON implements custom unmarshaling for StrapiData
// In Strapi 5, the attributes are embedded directly in the data object
func (s *StrapiData[T]) UnmarshalJSON(data []byte) error {
	fmt.Printf("[StrapiData UnmarshalJSON] Raw data: %s\n", string(data))

	// First, unmarshal the entire object into Attributes
	// This will capture all the content fields (Name, Slug, Logo, etc.)
	if err := json.Unmarshal(data, &s.Attributes); err != nil {
		fmt.Printf("[StrapiData UnmarshalJSON] Error unmarshaling attributes: %v\n", err)
		return err
	}

	fmt.Printf("[StrapiData UnmarshalJSON] Attributes after unmarshal: %+v\n", s.Attributes)

	// Then unmarshal metadata fields into a temporary struct
	var meta struct {
		ID          int        `json:"id"`
		DocumentID  string     `json:"documentId,omitempty"`
		CreatedAt   time.Time  `json:"createdAt,omitempty"`
		UpdatedAt   time.Time  `json:"updatedAt,omitempty"`
		PublishedAt *time.Time `json:"publishedAt,omitempty"`
		Locale      string     `json:"locale,omitempty"`
	}

	if err := json.Unmarshal(data, &meta); err != nil {
		return err
	}

	// Copy metadata to the struct
	s.ID = meta.ID
	s.DocumentID = meta.DocumentID
	s.CreatedAt = meta.CreatedAt
	s.UpdatedAt = meta.UpdatedAt
	s.PublishedAt = meta.PublishedAt
	s.Locale = meta.Locale

	return nil
}

// MarshalJSON implements custom marshaling for StrapiData
// This ensures the Attributes are embedded in the output alongside metadata
func (s StrapiData[T]) MarshalJSON() ([]byte, error) {
	// First marshal the attributes to get all content fields
	attributesJSON, err := json.Marshal(s.Attributes)
	if err != nil {
		return nil, err
	}

	// Parse it into a map
	var result map[string]interface{}
	if err := json.Unmarshal(attributesJSON, &result); err != nil {
		return nil, err
	}

	// Add metadata fields to the map
	result["id"] = s.ID
	if s.DocumentID != "" {
		result["documentId"] = s.DocumentID
	}
	if !s.CreatedAt.IsZero() {
		result["createdAt"] = s.CreatedAt
	}
	if !s.UpdatedAt.IsZero() {
		result["updatedAt"] = s.UpdatedAt
	}
	if s.PublishedAt != nil {
		result["publishedAt"] = s.PublishedAt
	}
	if s.Locale != "" {
		result["locale"] = s.Locale
	}

	// Marshal the combined map
	return json.Marshal(result)
}
