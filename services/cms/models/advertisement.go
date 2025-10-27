package models

// Advertisement represents an advertisement banner in the CMS
type Advertisement struct {
	Action string      `json:"Action"`
	Banner *MediaField `json:"Banner,omitempty"`
}

// AdvertisementResponse wraps a single Advertisement in Strapi's response format
type AdvertisementResponse = StrapiResponse[Advertisement]

// AdvertisementCollectionResponse wraps multiple Advertisements in Strapi's response format
type AdvertisementCollectionResponse = StrapiCollectionResponse[Advertisement]
