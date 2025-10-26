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
