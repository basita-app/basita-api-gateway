package models

// Common Strapi field types that can be reused across models

// MediaField represents a simplified Strapi 5 media field with only essential data
// Returns only: id (as documentId), width, height, url (prefixed), and formats
type MediaField struct {
	ID      string        `json:"id"`               // documentId from Strapi
	Width   int           `json:"width,omitempty"`  // Original width
	Height  int           `json:"height,omitempty"` // Original height
	URL     string        `json:"url"`              // Prefixed URL
	Formats *MediaFormats `json:"formats,omitempty"` // Available format sizes
}

// MediaCollectionField represents a collection of media fields
type MediaCollectionField []MediaField

// MediaFormats contains different sizes of media files
type MediaFormats struct {
	Thumbnail *MediaFormat `json:"thumbnail,omitempty"`
	Small     *MediaFormat `json:"small,omitempty"`
	Medium    *MediaFormat `json:"medium,omitempty"`
	Large     *MediaFormat `json:"large,omitempty"`
}

// MediaFormat represents a specific media format/size
type MediaFormat struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"` // Prefixed URL
}
