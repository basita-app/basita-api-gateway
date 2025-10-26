package models

// Common Strapi field types that can be reused across models

// MediaField represents a Strapi media field (single image, video, file)
type MediaField struct {
	Data *MediaData `json:"data,omitempty"`
}

// MediaCollectionField represents a Strapi media field (multiple files)
type MediaCollectionField struct {
	Data []MediaData `json:"data,omitempty"`
}

// MediaData contains the media file information
type MediaData struct {
	ID         int               `json:"id"`
	Attributes MediaAttributes   `json:"attributes"`
}

// MediaAttributes contains the media file attributes
type MediaAttributes struct {
	Name              string               `json:"name"`
	AlternativeText   string               `json:"alternativeText,omitempty"`
	Caption           string               `json:"caption,omitempty"`
	Width             int                  `json:"width,omitempty"`
	Height            int                  `json:"height,omitempty"`
	Formats           *MediaFormats        `json:"formats,omitempty"`
	Hash              string               `json:"hash"`
	Ext               string               `json:"ext"`
	Mime              string               `json:"mime"`
	Size              float64              `json:"size"`
	URL               string               `json:"url"`
	PreviewURL        string               `json:"previewUrl,omitempty"`
	Provider          string               `json:"provider"`
	ProviderMetadata  interface{}          `json:"provider_metadata,omitempty"`
}

// MediaFormats contains different sizes of media files
type MediaFormats struct {
	Thumbnail *MediaFormat `json:"thumbnail,omitempty"`
	Small     *MediaFormat `json:"small,omitempty"`
	Medium    *MediaFormat `json:"medium,omitempty"`
	Large     *MediaFormat `json:"large,omitempty"`
}

// MediaFormat represents a specific media format/size
type MediaFormat struct {
	Name   string  `json:"name"`
	Hash   string  `json:"hash"`
	Ext    string  `json:"ext"`
	Mime   string  `json:"mime"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Size   float64 `json:"size"`
	URL    string  `json:"url"`
}

// RelationField represents a Strapi relation field (single)
type RelationField[T any] struct {
	Data *StrapiData[T] `json:"data,omitempty"`
}

// RelationCollectionField represents a Strapi relation field (multiple)
type RelationCollectionField[T any] struct {
	Data []StrapiData[T] `json:"data,omitempty"`
}

// ComponentField represents a Strapi component field (single)
type ComponentField[T any] struct {
	Data T `json:"data,omitempty"`
}

// ComponentCollectionField represents a Strapi component field (repeatable)
type ComponentCollectionField[T any] struct {
	Data []T `json:"data,omitempty"`
}

// DynamicZoneItem represents an item in a Strapi dynamic zone
type DynamicZoneItem struct {
	ID        int                    `json:"id"`
	Component string                 `json:"__component"`
	Data      map[string]interface{} `json:"data,omitempty"`
}
