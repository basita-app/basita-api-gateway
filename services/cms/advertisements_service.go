package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"fmt"
)

// AdvertisementServiceGraphQL handles operations for advertisements using GraphQL
type AdvertisementServiceGraphQL struct {
	client *CMSClient
}

// NewAdvertisementServiceGraphQL creates a new GraphQL-based service for advertisements
func NewAdvertisementServiceGraphQL(client *CMSClient) *AdvertisementServiceGraphQL {
	return &AdvertisementServiceGraphQL{
		client: client,
	}
}

// parseMediaField converts Strapi media field to our MediaField model with URL prefixing
func (s *AdvertisementServiceGraphQL) parseMediaField(media *strapiMediaField) *models.MediaField {
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

// GetAll fetches all advertisements
func (s *AdvertisementServiceGraphQL) GetAll(ctx context.Context) (models.AdvertisementCollectionResponse, error) {
	data, err := s.client.ExecuteGraphQL(ctx, GetAdvertisementsQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch advertisements: %w", err)
	}

	var result struct {
		Advertisements []struct {
			DocumentID string            `json:"documentId"`
			Action     string            `json:"Action"`
			Banner     *strapiMediaField `json:"Banner"`
		} `json:"advertisements"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal advertisements: %w", err)
	}

	ads := make(models.AdvertisementCollectionResponse, len(result.Advertisements))
	for i, ad := range result.Advertisements {
		ads[i] = models.AdvertisementData{
			ID:     ad.DocumentID,
			Action: ad.Action,
			Banner: s.parseMediaField(ad.Banner),
		}
	}

	return ads, nil
}

// GetByID fetches a single advertisement by documentId
func (s *AdvertisementServiceGraphQL) GetByID(ctx context.Context, documentId string) (*models.AdvertisementResponse, error) {
	variables := map[string]interface{}{
		"documentId": documentId,
	}

	data, err := s.client.ExecuteGraphQL(ctx, GetAdvertisementByIDQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch advertisement: %w", err)
	}

	var result struct {
		Advertisement *struct {
			DocumentID string            `json:"documentId"`
			Action     string            `json:"Action"`
			Banner     *strapiMediaField `json:"Banner"`
		} `json:"advertisement"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal advertisement: %w", err)
	}

	if result.Advertisement == nil {
		return nil, fmt.Errorf("advertisement not found")
	}

	return &models.AdvertisementResponse{
		Data: &models.AdvertisementData{
			ID:     result.Advertisement.DocumentID,
			Action: result.Advertisement.Action,
			Banner: s.parseMediaField(result.Advertisement.Banner),
		},
	}, nil
}
