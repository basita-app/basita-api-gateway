package models

// CarModel represents a car model in the system
type CarModel struct {
	Name     string                     `json:"Name"`
	BodyType string                     `json:"BodyType"` // Sedan, SUV
	FuelType string                     `json:"FuelType,omitempty"` // Electric, Hybrid, x80, x90, x92, x95
	Slug     string                     `json:"Slug"`
	Brand    *RelationField[Brand]      `json:"Brand,omitempty"`
}

// CarModelResponse is a convenience type for car model API responses
type CarModelResponse = StrapiResponse[CarModel]

// CarModelCollectionResponse is a convenience type for car model collection API responses
type CarModelCollectionResponse = StrapiCollectionResponse[CarModel]
