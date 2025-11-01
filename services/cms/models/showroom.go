package models

// Showroom represents a showroom with its basic information
type Showroom struct {
	ID          string      `json:"id"`                    // documentId from Strapi
	Name        string      `json:"name"`                  // Showroom name
	Description string      `json:"description,omitempty"` // Showroom description
	IsVerified  bool        `json:"isVerified"`            // Verification status
	IsFeatured  bool        `json:"isFeatured"`            // Featured status
	Logo        *MediaField `json:"logo,omitempty"`        // Showroom logo
}

// DetailedShowroom represents a showroom with full details
type DetailedShowroom struct {
	ID              string        `json:"id"`                      // documentId from Strapi
	Name            string        `json:"name"`                    // Showroom name
	Description     string        `json:"description,omitempty"`   // Showroom description
	IsVerified      bool          `json:"isVerified"`              // Verification status
	IsFeatured      bool          `json:"isFeatured"`              // Featured status
	Logo            *MediaField   `json:"logo,omitempty"`          // Showroom logo
	Cover           *MediaField   `json:"cover,omitempty"`         // Showroom cover image
	OperatingHours  string        `json:"operatingHours,omitempty"` // Operating hours
	Location        *Location     `json:"location,omitempty"`      // Location details
	ContactInfo     *ContactInfo  `json:"contactInfo,omitempty"`   // Contact information
}

// Location represents showroom location details
type Location struct {
	Address     string       `json:"address,omitempty"`     // Street address
	Governorate *Governorate `json:"governorate,omitempty"` // Governorate/State
	City        *City        `json:"city,omitempty"`        // City
	Latitude    float64      `json:"latitude,omitempty"`    // GPS latitude
	Longitude   float64      `json:"longitude,omitempty"`   // GPS longitude
}

// Governorate represents a governorate/state
type Governorate struct {
	ID   string `json:"id"`   // documentId from Strapi
	Name string `json:"name"` // Governorate name
}

// City represents a city
type City struct {
	ID   string `json:"id"`   // documentId from Strapi
	Name string `json:"name"` // City name
}

// ContactInfo represents showroom contact information
type ContactInfo struct {
	Email      string `json:"email,omitempty"`      // Email address
	Phone      string `json:"phone,omitempty"`      // Phone number
	Facebook   string `json:"facebook,omitempty"`   // Facebook URL
	Instagram  string `json:"instagram,omitempty"`  // Instagram URL
	Tiktok     string `json:"tiktok,omitempty"`     // TikTok URL
	Whatsapp   string `json:"whatsapp,omitempty"`   // WhatsApp number
	X          string `json:"x,omitempty"`          // X (Twitter) URL
	Youtube    string `json:"youtube,omitempty"`    // YouTube URL
	WebsiteURL string `json:"websiteURL,omitempty"` // Website URL
}
