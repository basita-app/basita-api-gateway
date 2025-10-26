package models

// Showroom represents a showroom in the system
type Showroom struct {
	Name           string                `json:"name"`
	Description    string                `json:"description,omitempty"`
	Logo           *MediaField           `json:"logo,omitempty"`
	IsVerified     bool                  `json:"isVerified"`
	IsFeatured     bool                  `json:"isFeatured"`
	OperatingHours string                `json:"operatingHours,omitempty"`
	Location       *LocationComponent    `json:"location,omitempty"`
	ContactInfo    *ContactInfoComponent `json:"contactInfo,omitempty"`
	AvailableCars  []AvailableCarComponent `json:"availableCars,omitempty"`
}

// LocationComponent represents the location embedded component
type LocationComponent struct {
	Governorate *RelationField[Governorate] `json:"governorate,omitempty"`
	City        *RelationField[City]        `json:"city,omitempty"`
	Address     string                      `json:"address,omitempty"`
	Latitude    float64                     `json:"latitude,omitempty"`
	Longitude   float64                     `json:"longitude,omitempty"`
}

// ContactInfoComponent represents the contact info embedded component
type ContactInfoComponent struct {
	Phone      string `json:"phone,omitempty"`
	Whatsapp   string `json:"whatsapp,omitempty"`
	Email      string `json:"email,omitempty"`
	WebsiteURL string `json:"websiteURL,omitempty"`
	Tiktok     string `json:"tiktok,omitempty"`
	Youtube    string `json:"youtube,omitempty"`
	X          string `json:"x,omitempty"` // Twitter/X
	Instagram  string `json:"instagram,omitempty"`
	Facebook   string `json:"facebook,omitempty"`
}

// AvailableCarComponent represents a car available at the showroom
type AvailableCarComponent struct {
	Car   *RelationField[CarVariant] `json:"car,omitempty"`
	Price int                        `json:"price,omitempty"`
}

// ShowroomResponse is a convenience type for showroom API responses
type ShowroomResponse = StrapiResponse[Showroom]

// ShowroomCollectionResponse is a convenience type for showroom collection API responses
type ShowroomCollectionResponse = StrapiCollectionResponse[Showroom]
