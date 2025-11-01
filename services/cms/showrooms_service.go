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

// GetByID fetches a single showroom by documentId with full details
func (s *ShowroomServiceGraphQL) GetByID(ctx context.Context, documentId string) (*models.DetailedShowroom, error) {
	variables := map[string]interface{}{
		"documentId": documentId,
	}

	data, err := s.client.ExecuteGraphQL(ctx, GetShowroomByIDQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch showroom: %w", err)
	}

	var result struct {
		Showroom *struct {
			DocumentID     string            `json:"documentId"`
			Name           string            `json:"Name"`
			Description    string            `json:"Description"`
			IsVerified     bool              `json:"IsVerified"`
			IsFeatured     bool              `json:"IsFeatured"`
			Logo           *strapiMediaField `json:"Logo"`
			Cover          *strapiMediaField `json:"Cover"`
			OperatingHours string            `json:"OperatingHours"`
			Location       *struct {
				Address     string   `json:"Address"`
				Latitude    float64  `json:"Latitude"`
				Longitude   float64  `json:"Longitude"`
				Governorate *struct {
					DocumentID string `json:"documentId"`
					Name       string `json:"Name"`
				} `json:"governorate"`
				City *struct {
					DocumentID string `json:"documentId"`
					Name       string `json:"Name"`
				} `json:"city"`
			} `json:"Location"`
			ContactInfo *struct {
				Email      string `json:"Email"`
				Phone      string `json:"Phone"`
				Facebook   string `json:"Facebook"`
				Instagram  string `json:"Instagram"`
				Tiktok     string `json:"Tiktok"`
				Whatsapp   string `json:"Whatsapp"`
				X          string `json:"X"`
				Youtube    string `json:"Youtube"`
				WebsiteURL string `json:"WebsiteURL"`
			} `json:"ContactInfo"`
		} `json:"showroom"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal showroom: %w", err)
	}

	if result.Showroom == nil {
		return nil, fmt.Errorf("showroom not found")
	}

	showroom := &models.DetailedShowroom{
		ID:             result.Showroom.DocumentID,
		Name:           result.Showroom.Name,
		Description:    result.Showroom.Description,
		IsVerified:     result.Showroom.IsVerified,
		IsFeatured:     result.Showroom.IsFeatured,
		Logo:           s.parseMediaField(result.Showroom.Logo),
		Cover:          s.parseMediaField(result.Showroom.Cover),
		OperatingHours: result.Showroom.OperatingHours,
	}

	// Parse location
	if result.Showroom.Location != nil {
		location := &models.Location{
			Address:   result.Showroom.Location.Address,
			Latitude:  result.Showroom.Location.Latitude,
			Longitude: result.Showroom.Location.Longitude,
		}

		if result.Showroom.Location.Governorate != nil {
			location.Governorate = &models.Governorate{
				ID:   result.Showroom.Location.Governorate.DocumentID,
				Name: result.Showroom.Location.Governorate.Name,
			}
		}

		if result.Showroom.Location.City != nil {
			location.City = &models.City{
				ID:   result.Showroom.Location.City.DocumentID,
				Name: result.Showroom.Location.City.Name,
			}
		}

		showroom.Location = location
	}

	// Parse contact info
	if result.Showroom.ContactInfo != nil {
		showroom.ContactInfo = &models.ContactInfo{
			Email:      result.Showroom.ContactInfo.Email,
			Phone:      result.Showroom.ContactInfo.Phone,
			Facebook:   result.Showroom.ContactInfo.Facebook,
			Instagram:  result.Showroom.ContactInfo.Instagram,
			Tiktok:     result.Showroom.ContactInfo.Tiktok,
			Whatsapp:   result.Showroom.ContactInfo.Whatsapp,
			X:          result.Showroom.ContactInfo.X,
			Youtube:    result.Showroom.ContactInfo.Youtube,
			WebsiteURL: result.Showroom.ContactInfo.WebsiteURL,
		}
	}

	return showroom, nil
}

// GetCarVariantsByShowroomID fetches car variants available at a specific showroom
func (s *ShowroomServiceGraphQL) GetCarVariantsByShowroomID(ctx context.Context, showroomDocumentId string) ([]models.ShowroomVariant, error) {
	variables := map[string]interface{}{
		"showroomDocumentId": showroomDocumentId,
	}

	data, err := s.client.ExecuteGraphQL(ctx, GetCarVariantsByShowroomQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch variants: %w", err)
	}

	var result struct {
		CarVariants []struct {
			DocumentID      string             `json:"documentId"`
			DisplayName     string             `json:"DisplayName"`
			Images          []strapiMediaField `json:"Images"`
			ShowroomPricing []struct {
				Price           int `json:"Price"`
				MinDownPayment  int `json:"MinimuDownpayment"`
				MinInstallments int `json:"MinimumInstallements"`
			} `json:"ShowroomPricing"`
		} `json:"carVariants"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal variants: %w", err)
	}

	variants := make([]models.ShowroomVariant, 0, len(result.CarVariants))
	for _, variant := range result.CarVariants {
		showroomVariant := models.ShowroomVariant{
			ID:          variant.DocumentID,
			DisplayName: variant.DisplayName,
		}

		// Parse images
		if len(variant.Images) > 0 {
			images := make(models.MediaCollectionField, len(variant.Images))
			for i, img := range variant.Images {
				images[i] = *s.parseMediaField(&img)
			}
			showroomVariant.Images = &images
		}

		// Parse showroom pricing (should only be one since we filtered)
		if len(variant.ShowroomPricing) > 0 {
			pricing := variant.ShowroomPricing[0]
			showroomVariant.Price = pricing.Price
			showroomVariant.MinDownPayment = pricing.MinDownPayment
			showroomVariant.MinInstallments = pricing.MinInstallments
		}

		variants = append(variants, showroomVariant)
	}

	return variants, nil
}
