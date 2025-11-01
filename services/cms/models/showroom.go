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
