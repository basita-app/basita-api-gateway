package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"time"
)

// BrandService handles operations for the brands resource
type BrandService struct {
	client   Client
	endpoint string
	cacheTTL time.Duration
}

// NewBrandService creates a new service for brands
func NewBrandService(client Client) *BrandService {
	return &BrandService{
		client:   client,
		endpoint: "brands",
		cacheTTL: 30 * time.Minute, // Brands don't change often
	}
}

// GetAll fetches all brands with optional filters
func (s *BrandService) GetAll(ctx context.Context, opts models.CollectionOptions, useCache bool) (*models.BrandCollectionResponse, error) {
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

	var response models.BrandCollectionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal brands: " + err.Error(),
		}
	}

	return &response, nil
}

// GetByID fetches a single brand by ID
func (s *BrandService) GetByID(ctx context.Context, id string, opts models.ItemOptions, useCache bool) (*models.BrandResponse, error) {
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

	var response models.BrandResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal brand: " + err.Error(),
		}
	}

	return &response, nil
}

// InvalidateCache invalidates all cached brands
func (s *BrandService) InvalidateCache(ctx context.Context) error {
	return s.client.InvalidateCache(ctx, s.endpoint+"*")
}
