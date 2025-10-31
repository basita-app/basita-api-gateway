package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"fmt"
)

// CarModelServiceGraphQL handles operations for car models using GraphQL
type CarModelServiceGraphQL struct {
	client *CMSClient
}

// NewCarModelServiceGraphQL creates a new GraphQL-based service for car models
func NewCarModelServiceGraphQL(client *CMSClient) *CarModelServiceGraphQL {
	return &CarModelServiceGraphQL{
		client: client,
	}
}

// parseMediaField converts Strapi media field to our MediaField model with URL prefixing
func (s *CarModelServiceGraphQL) parseMediaField(media *strapiMediaField) *models.MediaField {
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

// GetByBrandID fetches car models for a specific brand with pricing
func (s *CarModelServiceGraphQL) GetByBrandID(ctx context.Context, brandDocumentID string) ([]models.SimpleCarModel, error) {
	variables := map[string]interface{}{
		"brandDocumentId": brandDocumentID,
	}

	data, err := s.client.ExecuteGraphQL(ctx, GetCarModelsByBrandQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch car models: %w", err)
	}

	var result struct {
		CarModels []struct {
			DocumentID string              `json:"documentId"`
			Name       string              `json:"Name"`
			Images     []strapiMediaField `json:"Images"`
		} `json:"carModels"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal car models: %w", err)
	}

	// For each car model, fetch variants to calculate prices
	carModels := make([]models.SimpleCarModel, len(result.CarModels))
	for i, model := range result.CarModels {
		// Get variants to calculate prices
		variants, err := s.getVariantsByModel(ctx, model.DocumentID)
		if err != nil {
			fmt.Printf("[WARN] Failed to fetch variants for model %s: %v\n", model.DocumentID, err)
		}

		// Calculate price range
		priceFrom := 0
		priceTo := 0
		if len(variants) > 0 {
			priceFrom = variants[0].Price
			priceTo = variants[0].Price
			for _, v := range variants {
				if v.Price < priceFrom {
					priceFrom = v.Price
				}
				if v.Price > priceTo {
					priceTo = v.Price
				}
			}
		}

		// Use first image as thumbnail
		var thumbnail *models.MediaField
		if len(model.Images) > 0 {
			thumbnail = s.parseMediaField(&model.Images[0])
		}

		carModels[i] = models.SimpleCarModel{
			ID:              model.DocumentID,
			Title:           model.Name,
			Thumbnail:       thumbnail,
			PriceFrom:       priceFrom,
			PriceTo:         priceTo,
			MarketPriceFrom: priceFrom + 20000,
			MarketPriceTo:   priceTo + 20000,
		}
	}

	return carModels, nil
}

// getVariantsByModel is a helper to fetch variants for price calculation
func (s *CarModelServiceGraphQL) getVariantsByModel(ctx context.Context, carModelDocumentID string) ([]struct {
	Price int `json:"Price"`
}, error) {
	query := `
		query GetVariantPrices($carModelDocumentId: ID!) {
			carVariants(filters: { car_model: { documentId: { eq: $carModelDocumentId } } }) {
				Price
			}
		}
	`

	variables := map[string]interface{}{
		"carModelDocumentId": carModelDocumentID,
	}

	data, err := s.client.ExecuteGraphQL(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var result struct {
		CarVariants []struct {
			Price int `json:"Price"`
		} `json:"carVariants"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result.CarVariants, nil
}

// GetDetailedByID fetches a detailed car model with all variants and showrooms
func (s *CarModelServiceGraphQL) GetDetailedByID(ctx context.Context, carModelDocumentID string) (*models.DetailedCarModel, error) {
	// Fetch car model info
	modelVars := map[string]interface{}{
		"documentId": carModelDocumentID,
	}

	modelData, err := s.client.ExecuteGraphQL(ctx, GetCarModelByIDQuery, modelVars)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch car model: %w", err)
	}

	var modelResult struct {
		CarModel *struct {
			DocumentID string              `json:"documentId"`
			Name       string              `json:"Name"`
			Images     []strapiMediaField `json:"Images"`
		} `json:"carModel"`
	}

	if err := json.Unmarshal(modelData, &modelResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal car model: %w", err)
	}

	if modelResult.CarModel == nil {
		return nil, fmt.Errorf("car model not found")
	}

	// Convert images
	images := make(models.MediaCollectionField, len(modelResult.CarModel.Images))
	for i, img := range modelResult.CarModel.Images {
		parsed := s.parseMediaField(&img)
		if parsed != nil {
			images[i] = *parsed
		}
	}

	// Fetch variants with showrooms
	variantVars := map[string]interface{}{
		"carModelDocumentId": carModelDocumentID,
	}

	variantData, err := s.client.ExecuteGraphQL(ctx, GetCarVariantsByModelQuery, variantVars)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch variants: %w", err)
	}

	var variantResult struct {
		CarVariants []struct {
			DocumentID          string `json:"documentId"`
			Name                string `json:"Name"`
			Price               int    `json:"Price"`
			Year                int    `json:"Year"`
			BrochureURL         string `json:"BrochureURL"`
			ReviewLink          string `json:"ReviewLink"`
			Warranty            string `json:"Warranty"`
			MinimumDownPayment  int    `json:"MinimumDownPaymet"`
			MinimumInstallments int    `json:"MinimumInstallments"`
			Specs               *struct {
				Motor string `json:"Motor"`
			} `json:"Specs"`
			ShowroomPricing []struct {
				Price    int `json:"Price"`
				Showroom *struct {
					DocumentID string            `json:"documentId"`
					Name       string            `json:"Name"`
					Logo       *strapiMediaField `json:"Logo"`
				} `json:"showroom"`
			} `json:"ShowroomPricing"`
		} `json:"carVariants"`
	}

	if err := json.Unmarshal(variantData, &variantResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal variants: %w", err)
	}

	// Build detailed response
	detailedModel := &models.DetailedCarModel{
		ID:              modelResult.CarModel.DocumentID,
		Title:           modelResult.CarModel.Name,
		Images:          &images,
		Variants:        make([]models.SimpleVariant, 0),
		Showrooms:       make([]models.SimpleShowroom, 0),
		Reviews:         make([]models.ReviewItem, 0),
		Catalogs:        make([]models.CatalogItem, 0),
	}

	if len(variantResult.CarVariants) == 0 {
		return detailedModel, nil
	}

	// Process variants
	priceFrom := variantResult.CarVariants[0].Price
	priceTo := variantResult.CarVariants[0].Price
	minDownPayment := 0
	minInstallments := 0
	warranty := ""

	showroomMap := make(map[string]*models.SimpleShowroom)
	reviewMap := make(map[string]*models.ReviewItem)
	catalogMap := make(map[string]*models.CatalogItem)
	reviewIDCounter := 1
	catalogIDCounter := 1

	for _, variant := range variantResult.CarVariants {
		// Update price range
		if variant.Price < priceFrom {
			priceFrom = variant.Price
		}
		if variant.Price > priceTo {
			priceTo = variant.Price
		}

		// Update min down payment and installments
		if variant.MinimumDownPayment > 0 {
			if minDownPayment == 0 || variant.MinimumDownPayment < minDownPayment {
				minDownPayment = variant.MinimumDownPayment
			}
		}
		if variant.MinimumInstallments > 0 {
			if minInstallments == 0 || variant.MinimumInstallments < minInstallments {
				minInstallments = variant.MinimumInstallments
			}
		}

		// Get warranty
		if warranty == "" && variant.Warranty != "" {
			warranty = variant.Warranty
		}

		// Add variant
		cc := ""
		if variant.Specs != nil {
			cc = variant.Specs.Motor
		}

		detailedModel.Variants = append(detailedModel.Variants, models.SimpleVariant{
			ID:    variant.DocumentID,
			Year:  variant.Year,
			Title: variant.Name,
			CC:    cc,
			Price: variant.Price,
		})

		// Process showrooms
		for _, pricing := range variant.ShowroomPricing {
			if pricing.Showroom != nil {
				showroomDocID := pricing.Showroom.DocumentID
				thumbnail := s.parseMediaField(pricing.Showroom.Logo)

				if existing, exists := showroomMap[showroomDocID]; exists {
					if pricing.Price < existing.Price {
						existing.Price = pricing.Price
						existing.MinDownPayment = variant.MinimumDownPayment
						existing.MinInstallments = variant.MinimumInstallments
					}
				} else {
					showroomMap[showroomDocID] = &models.SimpleShowroom{
						ID:              showroomDocID,
						Title:           pricing.Showroom.Name,
						Thumbnail:       thumbnail,
						Price:           pricing.Price,
						MinDownPayment:  variant.MinimumDownPayment,
						MinInstallments: variant.MinimumInstallments,
					}
				}
			}
		}

		// Process reviews
		if variant.ReviewLink != "" {
			if _, exists := reviewMap[variant.ReviewLink]; !exists {
				reviewMap[variant.ReviewLink] = &models.ReviewItem{
					ID:         reviewIDCounter,
					YoutubeURL: variant.ReviewLink,
				}
				reviewIDCounter++
			}
		}

		// Process catalogs
		if variant.BrochureURL != "" {
			if _, exists := catalogMap[variant.BrochureURL]; !exists {
				catalogMap[variant.BrochureURL] = &models.CatalogItem{
					ID:          catalogIDCounter,
					DownloadURL: variant.BrochureURL,
				}
				catalogIDCounter++
			}
		}
	}

	// Set calculated values
	detailedModel.PriceFrom = priceFrom
	detailedModel.PriceTo = priceTo
	detailedModel.MarketPriceFrom = priceFrom + 20000
	detailedModel.MarketPriceTo = priceTo + 20000
	detailedModel.MinDownPayment = minDownPayment
	detailedModel.MinInstallments = minInstallments
	detailedModel.Warranty = warranty

	// Convert maps to slices
	for _, showroom := range showroomMap {
		detailedModel.Showrooms = append(detailedModel.Showrooms, *showroom)
	}
	for _, review := range reviewMap {
		detailedModel.Reviews = append(detailedModel.Reviews, *review)
	}
	for _, catalog := range catalogMap {
		detailedModel.Catalogs = append(detailedModel.Catalogs, *catalog)
	}

	return detailedModel, nil
}
