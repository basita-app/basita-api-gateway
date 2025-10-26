package models

// Brand represents a car brand in the system
type Brand struct {
	Name string       `json:"name"`
	Slug string       `json:"slug"`
	Logo *MediaField  `json:"logo,omitempty"`
}

// BrandResponse is a convenience type for brand API responses
type BrandResponse = StrapiResponse[Brand]

// BrandCollectionResponse is a convenience type for brand collection API responses
type BrandCollectionResponse = StrapiCollectionResponse[Brand]
