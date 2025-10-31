package models

// SimpleBrand is a simplified brand response with only id, title, and thumbnail
type SimpleBrand struct {
	ID        string      `json:"id"`
	Title     string      `json:"title"`
	Thumbnail *MediaField `json:"thumbnail,omitempty"`
}
