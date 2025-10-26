package models

// Showroom represents a showroom in the system
type Showroom struct {
	Name           string                `json:"Name"`
	Description    string                `json:"Description,omitempty"`
	Logo           *MediaField           `json:"Logo,omitempty"`
	IsVerified     bool                  `json:"IsVerified"`
	IsFeatured     bool                  `json:"IsFeatured"`
	OperatingHours string                `json:"OperatingHours,omitempty"`
	Location       *LocationComponent    `json:"Location,omitempty"`
	ContactInfo    *ContactInfoComponent `json:"ContactInfo,omitempty"`
	AvailableCars  []AvailableCarComponent `json:"AvailableCars,omitempty"`
}

// LocationComponent represents the location embedded component
type LocationComponent struct {
	Governorate *RelationField[Governorate] `json:"Governorate,omitempty"`
	City        *RelationField[City]        `json:"City,omitempty"`
	Address     string                      `json:"Address,omitempty"`
	Latitude    float64                     `json:"Latitude,omitempty"`
	Longitude   float64                     `json:"Longitude,omitempty"`
}

// ContactInfoComponent represents the contact info embedded component
type ContactInfoComponent struct {
	Phone      string `json:"Phone,omitempty"`
	Whatsapp   string `json:"Whatsapp,omitempty"`
	Email      string `json:"Email,omitempty"`
	WebsiteURL string `json:"WebsiteURL,omitempty"`
	Tiktok     string `json:"Tiktok,omitempty"`
	Youtube    string `json:"Youtube,omitempty"`
	X          string `json:"X,omitempty"` // Twitter/X
	Instagram  string `json:"Instagram,omitempty"`
	Facebook   string `json:"Facebook,omitempty"`
}

// AvailableCarComponent represents a car available at the showroom
type AvailableCarComponent struct {
	Car   *RelationField[CarVariant] `json:"Car,omitempty"`
	Price int                        `json:"Price,omitempty"`
}

// ShowroomResponse is a convenience type for showroom API responses
type ShowroomResponse = StrapiResponse[Showroom]

// ShowroomCollectionResponse is a convenience type for showroom collection API responses
type ShowroomCollectionResponse = StrapiCollectionResponse[Showroom]
