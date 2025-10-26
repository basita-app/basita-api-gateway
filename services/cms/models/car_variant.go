package models

// CarVariant represents a car variant/trim in the system
type CarVariant struct {
	Name        string                         `json:"Name"`
	CarModel    *RelationField[CarModel]       `json:"CarModel,omitempty"`
	Price       int                            `json:"Price"`
	Year        int                            `json:"Year"`
	Images      *MediaCollectionField          `json:"Images,omitempty"`
	Specs       []SpecsComponent               `json:"Specs,omitempty"`
	Features    map[string]interface{}         `json:"Features,omitempty"`
	BrochureURL string                         `json:"BrochureURL,omitempty"`
	DisplayName string                         `json:"DisplayName,omitempty"`
}

// SpecsComponent represents the specs embedded component
type SpecsComponent struct {
	Motor                int     `json:"Motor,omitempty"`
	Speed                int     `json:"Speed,omitempty"`
	Transmission         string  `json:"Transmission,omitempty"`
	Horsepower           int     `json:"Horsepower,omitempty"`
	LiterPerKM           float64 `json:"LiterPerKM,omitempty"`
	MaxSpeed             int     `json:"MaxSpeed,omitempty"`
	Origin               string  `json:"Origin,omitempty"`
	AssembledIn          string  `json:"AssembledIn,omitempty"`
	Acceleration         float64 `json:"Acceleration,omitempty"`
	LengthInMM           int     `json:"LengthInMM,omitempty"`
	WidthInMM            int     `json:"WidthInMM,omitempty"`
	HeightInMM           int     `json:"HeightInMM,omitempty"`
	GroundClearanceInMM  int     `json:"GroundClearanceInMM,omitempty"`
	WheelBase            int     `json:"WheelBase,omitempty"`
}

// CarVariantResponse is a convenience type for car variant API responses
type CarVariantResponse = StrapiResponse[CarVariant]

// CarVariantCollectionResponse is a convenience type for car variant collection API responses
type CarVariantCollectionResponse = StrapiCollectionResponse[CarVariant]
