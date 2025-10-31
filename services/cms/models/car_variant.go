package models

// CarVariant represents a car variant/trim in the system
type CarVariant struct {
	Name        string                         `json:"Name"`
	CarModel    *RelationField[CarModel]       `json:"car_model,omitempty"`
	Price       int                            `json:"Price"`
	Year        int                            `json:"Year"`
	Images      *MediaCollectionField          `json:"Images,omitempty"`
	Specs       *SpecsComponent                `json:"Specs,omitempty"`
	Features    []string                       `json:"Features,omitempty"`
	BrochureURL string                         `json:"BrochureURL,omitempty"`
	DisplayName string                         `json:"DisplayName,omitempty"`
	ShowroomPricing []ShowroomPricingComponent `json:"ShowroomPricing,omitempty"`
}

// ShowroomPricingComponent represents showroom pricing
type ShowroomPricingComponent struct {
	ID    int `json:"id,omitempty"`
	Price int `json:"Price"`
}

// SpecsComponent represents the specs embedded component
type SpecsComponent struct {
	ID                  int     `json:"id,omitempty"`
	Motor               string  `json:"Motor,omitempty"`
	Speed               int     `json:"Speed,omitempty"`
	Transmission        string  `json:"Transmission,omitempty"`
	Horsepower          int     `json:"Horsepower,omitempty"`
	LiterPerKM          float64 `json:"LiterPerKM,omitempty"`
	MaxSpeed            int     `json:"MaxSpeed,omitempty"`
	Origin              string  `json:"Origin,omitempty"`
	AssembledIn         string  `json:"AssembledIn,omitempty"`
	Acceleration        float64 `json:"Acceleration,omitempty"`
	LengthInMM          int     `json:"LengthInMM,omitempty"`
	WidthInMM           int     `json:"WidthInMM,omitempty"`
	HeightInMM          int     `json:"HeightInMM,omitempty"`
	GroundClearanceInMM int     `json:"GroundClearanceInMM,omitempty"`
	WheelBase           int     `json:"WheelBase,omitempty"`
	TrunkSize           int     `json:"TrunkSize,omitempty"`
	Seats               int     `json:"Seats,omitempty"`
	TractionType        string  `json:"TractionType,omitempty"`
}

// CarVariantResponse is a convenience type for car variant API responses
type CarVariantResponse = StrapiResponse[CarVariant]

// CarVariantCollectionResponse is a convenience type for car variant collection API responses
type CarVariantCollectionResponse = StrapiCollectionResponse[CarVariant]
