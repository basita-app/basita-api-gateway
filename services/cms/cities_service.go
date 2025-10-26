package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"time"
)

// CityService handles operations for the cities resource
type CityService struct {
	client   Client
	endpoint string
	cacheTTL time.Duration
}

// NewCityService creates a new service for cities
func NewCityService(client Client) *CityService {
	return &CityService{
		client:   client,
		endpoint: "cities",
		cacheTTL: 60 * time.Minute, // Cities rarely change
	}
}

// GetAll fetches all cities with optional filters
func (s *CityService) GetAll(ctx context.Context, opts models.CollectionOptions, useCache bool) (*models.CityCollectionResponse, error) {
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

	var response models.CityCollectionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal cities: " + err.Error(),
		}
	}

	return &response, nil
}

// GetByID fetches a single city by ID
func (s *CityService) GetByID(ctx context.Context, id string, opts models.ItemOptions, useCache bool) (*models.CityResponse, error) {
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

	var response models.CityResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal city: " + err.Error(),
		}
	}

	return &response, nil
}

// InvalidateCache invalidates all cached cities
func (s *CityService) InvalidateCache(ctx context.Context) error {
	return s.client.InvalidateCache(ctx, s.endpoint+"*")
}
