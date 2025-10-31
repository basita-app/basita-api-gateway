package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"fmt"
)

// BrandServiceGraphQL handles operations for brands using GraphQL
type BrandServiceGraphQL struct {
	client *CMSClient
}

// NewBrandServiceGraphQL creates a new GraphQL-based service for brands
func NewBrandServiceGraphQL(client *CMSClient) *BrandServiceGraphQL {
	return &BrandServiceGraphQL{
		client: client,
	}
}

// parseMediaField converts Strapi media field to our MediaField model with URL prefixing
func (s *BrandServiceGraphQL) parseMediaField(media *strapiMediaField) *models.MediaField {
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

// GetSimplified fetches all brands with simplified response
func (s *BrandServiceGraphQL) GetSimplified(ctx context.Context) ([]models.SimpleBrand, error) {
	data, err := s.client.ExecuteGraphQL(ctx, GetBrandsQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch brands: %w", err)
	}

	var result struct {
		Brands []struct {
			DocumentID string            `json:"documentId"`
			Name       string            `json:"Name"`
			Logo       *strapiMediaField `json:"Logo"`
		} `json:"brands"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal brands: %w", err)
	}

	brands := make([]models.SimpleBrand, len(result.Brands))
	for i, brand := range result.Brands {
		brands[i] = models.SimpleBrand{
			ID:        brand.DocumentID,
			Title:     brand.Name,
			Thumbnail: s.parseMediaField(brand.Logo),
		}
	}

	return brands, nil
}

// GetByID fetches a single brand by documentId
func (s *BrandServiceGraphQL) GetByID(ctx context.Context, documentId string) (*models.SimpleBrand, error) {
	variables := map[string]interface{}{
		"documentId": documentId,
	}

	data, err := s.client.ExecuteGraphQL(ctx, GetBrandByIDQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch brand: %w", err)
	}

	var result struct {
		Brand *struct {
			DocumentID string            `json:"documentId"`
			Name       string            `json:"Name"`
			Slug       string            `json:"Slug"`
			Logo       *strapiMediaField `json:"Logo"`
		} `json:"brand"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal brand: %w", err)
	}

	if result.Brand == nil {
		return nil, fmt.Errorf("brand not found")
	}

	return &models.SimpleBrand{
		ID:        result.Brand.DocumentID,
		Title:     result.Brand.Name,
		Thumbnail: s.parseMediaField(result.Brand.Logo),
	}, nil
}
