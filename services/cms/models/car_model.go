package models

// CarModel represents a car model in the system
type CarModel struct {
	Name     string                     `json:"Name"`
	BodyType string                     `json:"BodyType"` // Sedan, SUV
	FuelType string                     `json:"FuelType,omitempty"` // Electric, Hybrid, x80, x90, x92, x95
	Slug     string                     `json:"Slug"`
	Brand    *RelationField[Brand]      `json:"brand,omitempty"`
	Image    *MediaField                `json:"Image,omitempty"`
}

// CarModelResponse is a convenience type for car model API responses
type CarModelResponse = StrapiResponse[CarModel]

// CarModelCollectionResponse is a convenience type for car model collection API responses
type CarModelCollectionResponse = StrapiCollectionResponse[CarModel]

// SimpleCarModel is a simplified car model response with pricing
type SimpleCarModel struct {
	ID                int         `json:"id"`
	Title             string      `json:"title"`
	Thumbnail         *MediaField `json:"thumbnail,omitempty"`
	PriceFrom         int         `json:"pricefrom"`
	PriceTo           int         `json:"priceto"`
	MarketPriceFrom   int         `json:"marketpricefrom"`
	MarketPriceTo     int         `json:"marketpriceto"`
}
