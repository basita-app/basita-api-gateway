package models

// City represents a city in the system
type City struct {
	Name        string                        `json:"Name"`
	Governorate *RelationField[Governorate]   `json:"Governorate,omitempty"`
}

// CityResponse is a convenience type for city API responses
type CityResponse = StrapiResponse[City]

// CityCollectionResponse is a convenience type for city collection API responses
type CityCollectionResponse = StrapiCollectionResponse[City]
