package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// CarModelService handles operations for the car-models resource
type CarModelService struct {
	client   Client
	endpoint string
	cacheTTL time.Duration
}

// NewCarModelService creates a new service for car models
func NewCarModelService(client Client) *CarModelService {
	return &CarModelService{
		client:   client,
		endpoint: "car-models",
		cacheTTL: 20 * time.Minute,
	}
}

// GetAll fetches all car models with optional filters
func (s *CarModelService) GetAll(ctx context.Context, opts models.CollectionOptions, useCache bool) (*models.CarModelCollectionResponse, error) {
	var cacheOpts *models.CacheOptions
	if useCache {
		cacheOpts = &models.CacheOptions{
			Enabled: true,
			TTL:     s.cacheTTL,
		}
	}

	data, err := s.client.GetCollection(ctx, s.endpoint, opts, cacheOpts)
	if err != nil {
		return nil, err
	}

	var response models.CarModelCollectionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal car models: " + err.Error(),
		}
	}

	return &response, nil
}

// GetByID fetches a single car model by ID
func (s *CarModelService) GetByID(ctx context.Context, id string, opts models.ItemOptions, useCache bool) (*models.CarModelResponse, error) {
	var cacheOpts *models.CacheOptions
	if useCache {
		cacheOpts = &models.CacheOptions{
			Enabled: true,
			TTL:     s.cacheTTL,
		}
	}

	data, err := s.client.GetItem(ctx, s.endpoint, id, opts, cacheOpts)
	if err != nil {
		return nil, err
	}

	var response models.CarModelResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal car model: " + err.Error(),
		}
	}

	return &response, nil
}

// InvalidateCache invalidates all cached car models
func (s *CarModelService) InvalidateCache(ctx context.Context) error {
	return s.client.InvalidateCache(ctx, s.endpoint+"*")
}

// GetByBrandID fetches all car models for a specific brand with pricing
func (s *CarModelService) GetByBrandID(ctx context.Context, brandID string, useCache bool) ([]models.SimpleCarModel, error) {
	var cacheOpts *models.CacheOptions
	if useCache {
		cacheOpts = &models.CacheOptions{
			Enabled: true,
			TTL:     s.cacheTTL,
		}
	}

	// Fetch car models filtered by brand ID
	opts := models.CollectionOptions{
		PageSize: 1000,
		Page:     1,
		Populate: "*", // Populate brand and image relations
		Filters: map[string]string{
			"filters[brand][id][$eq]": brandID,
		},
	}

	data, err := s.client.GetCollection(ctx, s.endpoint, opts, cacheOpts)
	if err != nil {
		return nil, err
	}

	var response models.CarModelCollectionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal car models: " + err.Error(),
		}
	}

	// Transform to simplified response with pricing
	// Initialize CarVariantService to fetch prices
	variantService := NewCarVariantService(s.client)
	
	simplifiedModels := make([]models.SimpleCarModel, len(response.Data))
	for i, modelData := range response.Data {
		// Fetch variants for this car model to calculate min/max prices
		variantOpts := models.CollectionOptions{
			PageSize: 1000,
			Page:     1,
			Populate: "*",
			Filters: map[string]string{
				"filters[car_model][id][$eq]": fmt.Sprintf("%d", modelData.ID),
			},
		}
		
		priceFrom := 0
		priceTo := 0

		// Get variants and calculate min/max prices
		variantData, err := variantService.getAllVariantsForModel(ctx, variantOpts, cacheOpts)
		if err == nil && len(variantData) > 0 {
			priceFrom = variantData[0].Attributes.Price
			priceTo = variantData[0].Attributes.Price
			
			for _, variant := range variantData {
				if variant.Attributes.Price < priceFrom {
					priceFrom = variant.Attributes.Price
				}
				if variant.Attributes.Price > priceTo {
					priceTo = variant.Attributes.Price
				}
			}
		}

		// Extract thumbnail image from the first image if available
		var thumbnail *models.MediaField
		if modelData.Attributes.Images != nil && len(*modelData.Attributes.Images) > 0 {
			firstImage := (*modelData.Attributes.Images)[0]
			if firstImage.Formats != nil && firstImage.Formats.Thumbnail != nil {
				// Convert MediaFormat to MediaField
				mediaFormat := firstImage.Formats.Thumbnail
				thumbnail = &models.MediaField{
					Name:   mediaFormat.Name,
					Hash:   mediaFormat.Hash,
					Ext:    mediaFormat.Ext,
					Mime:   mediaFormat.Mime,
					Width:  mediaFormat.Width,
					Height: mediaFormat.Height,
					Size:   mediaFormat.Size,
					URL:    mediaFormat.URL,
				}
			}
		}

		simplifiedModels[i] = models.SimpleCarModel{
			ID:              modelData.ID,
			Title:           modelData.Attributes.Name,
			Thumbnail:       thumbnail,
			PriceFrom:       priceFrom,
			PriceTo:         priceTo,
			MarketPriceFrom: 0, // Will be set in handler
			MarketPriceTo:   0, // Will be set in handler
		}
	}

	return simplifiedModels, nil
}

// getAllVariantsForModel is an internal method to fetch all variants without response wrapping
func (s *CarVariantService) getAllVariantsForModel(ctx context.Context, opts models.CollectionOptions, cacheOpts *models.CacheOptions) ([]models.StrapiData[models.CarVariant], error) {
	data, err := s.client.GetCollection(ctx, s.endpoint, opts, cacheOpts)
	if err != nil {
		return nil, err
	}

	var response models.CarVariantCollectionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal car variants: " + err.Error(),
		}
	}

	return response.Data, nil
}
