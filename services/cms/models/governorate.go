package models

// Governorate represents a governorate in the system
type Governorate struct {
	Name string `json:"Name"`
}

// GovernorateResponse is a convenience type for governorate API responses
type GovernorateResponse = StrapiResponse[Governorate]

// GovernorateCollectionResponse is a convenience type for governorate collection API responses
type GovernorateCollectionResponse = StrapiCollectionResponse[Governorate]
