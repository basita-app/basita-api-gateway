package models

// CarModel represents a car model in the system
type CarModel struct {
	Name     string                     `json:"name"`
	BodyType string                     `json:"bodyType"` // Sedan, SUV
	FuelType string                     `json:"fuelType,omitempty"` // Electric, Hybrid, x80, x90, x92, x95
	Slug     string                     `json:"slug"`
	Brand    *RelationField[Brand]      `json:"brand,omitempty"`
}

// CarModelResponse is a convenience type for car model API responses
type CarModelResponse = StrapiResponse[CarModel]

// CarModelCollectionResponse is a convenience type for car model collection API responses
type CarModelCollectionResponse = StrapiCollectionResponse[CarModel]
