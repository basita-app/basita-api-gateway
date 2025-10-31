package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
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

	// Fetch car models filtered by brand documentId
	opts := models.CollectionOptions{
		PageSize: 1000,
		Page:     1,
		Populate: "*", // Populate brand and image relations
		Filters: map[string]string{
			"filters[brand][documentId][$eq]": brandID,
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
			Populate: "deep", // Use deep population to get all nested relations
			Filters: map[string]string{
				"filters[car_model][documentId][$eq]": modelData.DocumentID,
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
			ID:              modelData.DocumentID,
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

// GetDetailedByID fetches a detailed car model by ID with all variants and showrooms
func (s *CarModelService) GetDetailedByID(ctx context.Context, id string, useCache bool) (*models.DetailedCarModel, error) {
	var cacheOpts *models.CacheOptions
	if useCache {
		cacheOpts = &models.CacheOptions{
			Enabled: true,
			TTL:     s.cacheTTL,
		}
	}

	// Fetch the car model
	opts := models.ItemOptions{
		Populate: "*",
	}

	data, err := s.client.GetItem(ctx, s.endpoint, id, opts, cacheOpts)
	if err != nil {
		return nil, err
	}

	var carModelResponse models.CarModelResponse
	if err := json.Unmarshal(data, &carModelResponse); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal car model: " + err.Error(),
		}
	}

	// Fetch variants for this car model
	variantService := NewCarVariantService(s.client)
	variantOpts := models.CollectionOptions{
		PageSize: 1000,
		Page:     1,
		Populate: "deep", // Use deep population to get nested ShowroomPricing.Showroom relations
		Filters: map[string]string{
			"filters[car_model][documentId][$eq]": id,
		},
	}

	variantData, err := variantService.getAllVariantsForModel(ctx, variantOpts, cacheOpts)
	if err != nil {
		return nil, err
	}

	// Initialize response
	detailedModel := &models.DetailedCarModel{
		ID:              carModelResponse.Data.DocumentID,
		Title:           carModelResponse.Data.Attributes.Name,
		Images:          carModelResponse.Data.Attributes.Images,
		Variants:        make([]models.SimpleVariant, 0),
		Showrooms:       make([]models.SimpleShowroom, 0),
		Reviews:         make([]models.ReviewItem, 0),
		Catalogs:        make([]models.CatalogItem, 0),
	}

	if len(variantData) == 0 {
		return detailedModel, nil
	}

	// Process variants to calculate prices and collect data
	priceFrom := variantData[0].Attributes.Price
	priceTo := variantData[0].Attributes.Price
	minDownPayment := 0
	minInstallments := 0
	warranty := ""

	// Maps to track unique showrooms and reviews
	showroomMap := make(map[string]*models.SimpleShowroom)
	reviewMap := make(map[string]*models.ReviewItem)
	catalogMap := make(map[string]*models.CatalogItem)
	reviewIDCounter := 1
	catalogIDCounter := 1

	for _, variant := range variantData {
		// Update price range
		if variant.Attributes.Price < priceFrom {
			priceFrom = variant.Attributes.Price
		}
		if variant.Attributes.Price > priceTo {
			priceTo = variant.Attributes.Price
		}

		// Update min down payment and installments
		if variant.Attributes.MinDownPayment > 0 {
			if minDownPayment == 0 || variant.Attributes.MinDownPayment < minDownPayment {
				minDownPayment = variant.Attributes.MinDownPayment
			}
		}
		if variant.Attributes.MinInstallments > 0 {
			if minInstallments == 0 || variant.Attributes.MinInstallments < minInstallments {
				minInstallments = variant.Attributes.MinInstallments
			}
		}

		// Get warranty (use first non-empty warranty)
		if warranty == "" && variant.Attributes.Warranty != "" {
			warranty = variant.Attributes.Warranty
		}

		// Add variant to list
		cc := ""
		if variant.Attributes.Specs != nil {
			cc = variant.Attributes.Specs.Motor
		}

		detailedModel.Variants = append(detailedModel.Variants, models.SimpleVariant{
			ID:    variant.DocumentID,
			Year:  variant.Attributes.Year,
			Title: variant.Attributes.Name,
			CC:    cc,
			Price: variant.Attributes.Price,
		})

		// Process showrooms from ShowroomPricing
		if variant.Attributes.ShowroomPricing != nil {
			for _, pricing := range variant.Attributes.ShowroomPricing {
				if pricing.Showroom != nil && pricing.Showroom.Data != nil && pricing.Showroom.Data.DocumentID != "" {
					showroomDocID := pricing.Showroom.Data.DocumentID
					showroom := pricing.Showroom.Data.Attributes

					// Extract thumbnail from showroom logo
					var thumbnail *models.MediaField
					if showroom.Logo != nil {
						thumbnail = showroom.Logo
					}

					// Check if showroom already exists, keep the one with minimum price
					if existing, exists := showroomMap[showroomDocID]; exists {
						if pricing.Price < existing.Price {
							existing.Price = pricing.Price
							existing.MinDownPayment = pricing.MinDownPayment
							existing.MinInstallments = pricing.MinInstallments
						}
					} else {
						showroomMap[showroomDocID] = &models.SimpleShowroom{
							ID:              showroomDocID,
							Title:           showroom.Name,
							Thumbnail:       thumbnail,
							Price:           pricing.Price,
							MinDownPayment:  pricing.MinDownPayment,
							MinInstallments: pricing.MinInstallments,
						}
					}
				}
			}
		}

		// Process review links
		if variant.Attributes.ReviewLink != "" {
			if _, exists := reviewMap[variant.Attributes.ReviewLink]; !exists {
				reviewMap[variant.Attributes.ReviewLink] = &models.ReviewItem{
					ID:         reviewIDCounter,
					YoutubeURL: variant.Attributes.ReviewLink,
				}
				reviewIDCounter++
			}
		}

		// Process catalogs/brochures
		if variant.Attributes.BrochureURL != "" {
			if _, exists := catalogMap[variant.Attributes.BrochureURL]; !exists {
				catalogMap[variant.Attributes.BrochureURL] = &models.CatalogItem{
					ID:          catalogIDCounter,
					DownloadURL: variant.Attributes.BrochureURL,
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
