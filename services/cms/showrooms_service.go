package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"fmt"
)

// ShowroomServiceGraphQL handles operations for showrooms using GraphQL
type ShowroomServiceGraphQL struct {
	client *CMSClient
}

// NewShowroomServiceGraphQL creates a new GraphQL-based service for showrooms
func NewShowroomServiceGraphQL(client *CMSClient) *ShowroomServiceGraphQL {
	return &ShowroomServiceGraphQL{
		client: client,
	}
}

// parseMediaField converts Strapi media field to our MediaField model with URL prefixing
func (s *ShowroomServiceGraphQL) parseMediaField(media *strapiMediaField) *models.MediaField {
	if media == nil {
		return nil
	}

	field := &models.MediaField{
		ID:     media.DocumentID,
		Width:  media.Width,
		Height: media.Height,
		URL:    s.client.PrefixMediaURL(media.URL),
	}

	// Parse formats if available
	if media.Formats != nil {
		formats := &models.MediaFormats{}

		// Helper to parse individual format
		parseFormat := func(formatData interface{}) *models.MediaFormat {
			if formatData == nil {
				return nil
			}
			formatMap, ok := formatData.(map[string]interface{})
			if !ok {
				return nil
			}
			url, _ := formatMap["url"].(string)
			width, _ := formatMap["width"].(float64)
			height, _ := formatMap["height"].(float64)
			return &models.MediaFormat{
				Width:  int(width),
				Height: int(height),
				URL:    s.client.PrefixMediaURL(url),
			}
		}

		if thumbnail, ok := media.Formats["thumbnail"]; ok {
			formats.Thumbnail = parseFormat(thumbnail)
		}
		if small, ok := media.Formats["small"]; ok {
			formats.Small = parseFormat(small)
		}
		if medium, ok := media.Formats["medium"]; ok {
			formats.Medium = parseFormat(medium)
		}
		if large, ok := media.Formats["large"]; ok {
			formats.Large = parseFormat(large)
		}

		field.Formats = formats
	}

	return field
}

// GetAll fetches all showrooms
func (s *ShowroomServiceGraphQL) GetAll(ctx context.Context) ([]models.Showroom, error) {
	data, err := s.client.ExecuteGraphQL(ctx, GetShowroomsQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch showrooms: %w", err)
	}

	var result struct {
		Showrooms []struct {
			DocumentID  string            `json:"documentId"`
			Name        string            `json:"Name"`
			Description string            `json:"Description"`
			IsVerified  bool              `json:"IsVerified"`
			IsFeatured  bool              `json:"IsFeatured"`
			Logo        *strapiMediaField `json:"Logo"`
		} `json:"showrooms"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal showrooms: %w", err)
	}

	showrooms := make([]models.Showroom, len(result.Showrooms))
	for i, showroom := range result.Showrooms {
		showrooms[i] = models.Showroom{
			ID:          showroom.DocumentID,
			Name:        showroom.Name,
			Description: showroom.Description,
			IsVerified:  showroom.IsVerified,
			IsFeatured:  showroom.IsFeatured,
			Logo:        s.parseMediaField(showroom.Logo),
		}
	}

	return showrooms, nil
}
