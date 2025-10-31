package models

// AdvertisementData represents a simplified advertisement with only essential fields
type AdvertisementData struct {
	ID     string      `json:"id"`
	Action string      `json:"action"`
	Banner *MediaField `json:"banner,omitempty"`
}

// AdvertisementResponse represents a single advertisement response with only id, action, and banner
type AdvertisementResponse struct {
	Data *AdvertisementData `json:"data"`
}

// AdvertisementCollectionResponse represents all advertisements without pagination
type AdvertisementCollectionResponse []AdvertisementData
