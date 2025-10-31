package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"time"
)

// AdvertisementService handles operations for the advertisements resource
type AdvertisementService struct {
	client   Client
	endpoint string
	cacheTTL time.Duration
}

// NewAdvertisementService creates a new service for advertisements
func NewAdvertisementService(client Client) *AdvertisementService {
	return &AdvertisementService{
		client:   client,
		endpoint: "advertisements",
		cacheTTL: 15 * time.Minute, // Shorter cache for ads that may change frequently
	}
}

// GetAll fetches all advertisements without pagination
func (s *AdvertisementService) GetAll(ctx context.Context, opts models.CollectionOptions, useCache bool) (models.AdvertisementCollectionResponse, error) {
	var cacheOpts *models.CacheOptions
	if useCache {
		cacheOpts = &models.CacheOptions{
			Enabled: true,
			TTL:     s.cacheTTL,
		}
	}

	// Override pagination to fetch all advertisements at once
	opts.PageSize = 1000 // Set a high limit to get all advertisements

	data, err := s.client.GetCollection(ctx, s.endpoint, opts, cacheOpts)
	if err != nil {
		return nil, err
	}

	// Unmarshal into full Strapi response
	var fullResponse models.StrapiCollectionResponse[models.Advertisement]
	if err := json.Unmarshal(data, &fullResponse); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal advertisements: " + err.Error(),
		}
	}

	// Transform to simplified response with only id, action, and banner (no pagination)
	simplifiedData := make(models.AdvertisementCollectionResponse, len(fullResponse.Data))
	for i, item := range fullResponse.Data {
		simplifiedData[i] = models.AdvertisementData{
			ID:     item.DocumentID,
			Action: item.Attributes.Action,
			Banner: item.Attributes.Banner,
		}
	}

	return simplifiedData, nil
}

// GetByID fetches a single advertisement by ID
func (s *AdvertisementService) GetByID(ctx context.Context, id string, opts models.ItemOptions, useCache bool) (*models.AdvertisementResponse, error) {
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

	// Unmarshal into full Strapi response
	var fullResponse models.StrapiResponse[models.Advertisement]
	if err := json.Unmarshal(data, &fullResponse); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal advertisement: " + err.Error(),
		}
	}

	// Transform to simplified response with only id, action, and banner
	var simplifiedData *models.AdvertisementData
	if fullResponse.Data != nil {
		simplifiedData = &models.AdvertisementData{
			ID:     fullResponse.Data.DocumentID,
			Action: fullResponse.Data.Attributes.Action,
			Banner: fullResponse.Data.Attributes.Banner,
		}
	}

	return &models.AdvertisementResponse{
		Data: simplifiedData,
	}, nil
}

// InvalidateCache invalidates all cached advertisements
func (s *AdvertisementService) InvalidateCache(ctx context.Context) error {
	return s.client.InvalidateCache(ctx, s.endpoint+"*")
}
