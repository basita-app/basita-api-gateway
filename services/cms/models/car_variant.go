package models

// CarVariant represents a car variant/trim in the system
type CarVariant struct {
	Name        string                         `json:"name"`
	CarModel    *RelationField[CarModel]       `json:"carModel,omitempty"`
	Price       int                            `json:"price"`
	Year        int                            `json:"year"`
	Images      *MediaCollectionField          `json:"images,omitempty"`
	Specs       []SpecsComponent               `json:"specs,omitempty"`
	Features    map[string]interface{}         `json:"features,omitempty"`
	BrochureURL string                         `json:"brochureURL,omitempty"`
	DisplayName string                         `json:"displayName,omitempty"`
}

// SpecsComponent represents the specs embedded component
type SpecsComponent struct {
	Motor                int     `json:"motor,omitempty"`
	Speed                int     `json:"speed,omitempty"`
	Transmission         string  `json:"transmission,omitempty"`
	Horsepower           int     `json:"horsepower,omitempty"`
	LiterPerKM           float64 `json:"literPerKM,omitempty"`
	MaxSpeed             int     `json:"maxSpeed,omitempty"`
	Origin               string  `json:"origin,omitempty"`
	AssembledIn          string  `json:"assembledIn,omitempty"`
	Acceleration         float64 `json:"acceleration,omitempty"`
	LengthInMM           int     `json:"lengthInMM,omitempty"`
	WidthInMM            int     `json:"widthInMM,omitempty"`
	HeightInMM           int     `json:"heightInMM,omitempty"`
	GroundClearanceInMM  int     `json:"groundClearanceInMM,omitempty"`
	WheelBase            int     `json:"wheelBase,omitempty"`
}

// CarVariantResponse is a convenience type for car variant API responses
type CarVariantResponse = StrapiResponse[CarVariant]

// CarVariantCollectionResponse is a convenience type for car variant collection API responses
type CarVariantCollectionResponse = StrapiCollectionResponse[CarVariant]
