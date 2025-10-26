package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"time"
)

// CarVariantService handles operations for the car-variants resource
type CarVariantService struct {
	client   Client
	endpoint string
	cacheTTL time.Duration
}

// NewCarVariantService creates a new service for car variants
func NewCarVariantService(client Client) *CarVariantService {
	return &CarVariantService{
		client:   client,
		endpoint: "car-variants",
		cacheTTL: 15 * time.Minute,
	}
}

// GetAll fetches all car variants with optional filters
func (s *CarVariantService) GetAll(ctx context.Context, opts models.CollectionOptions, useCache bool) (*models.CarVariantCollectionResponse, error) {
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

	var response models.CarVariantCollectionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal car variants: " + err.Error(),
		}
	}

	return &response, nil
}

// GetByID fetches a single car variant by ID
func (s *CarVariantService) GetByID(ctx context.Context, id string, opts models.ItemOptions, useCache bool) (*models.CarVariantResponse, error) {
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

	var response models.CarVariantResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal car variant: " + err.Error(),
		}
	}

	return &response, nil
}

// InvalidateCache invalidates all cached car variants
func (s *CarVariantService) InvalidateCache(ctx context.Context) error {
	return s.client.InvalidateCache(ctx, s.endpoint+"*")
}
