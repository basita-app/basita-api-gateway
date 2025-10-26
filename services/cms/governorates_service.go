package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"time"
)

// GovernorateService handles operations for the governorates resource
type GovernorateService struct {
	client   Client
	endpoint string
	cacheTTL time.Duration
}

// NewGovernorateService creates a new service for governorates
func NewGovernorateService(client Client) *GovernorateService {
	return &GovernorateService{
		client:   client,
		endpoint: "governorates",
		cacheTTL: 60 * time.Minute, // Governorates rarely change
	}
}

// GetAll fetches all governorates with optional filters
func (s *GovernorateService) GetAll(ctx context.Context, opts models.CollectionOptions, useCache bool) (*models.GovernorateCollectionResponse, error) {
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

	var response models.GovernorateCollectionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal governorates: " + err.Error(),
		}
	}

	return &response, nil
}

// GetByID fetches a single governorate by ID
func (s *GovernorateService) GetByID(ctx context.Context, id string, opts models.ItemOptions, useCache bool) (*models.GovernorateResponse, error) {
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

	var response models.GovernorateResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal governorate: " + err.Error(),
		}
	}

	return &response, nil
}

// InvalidateCache invalidates all cached governorates
func (s *GovernorateService) InvalidateCache(ctx context.Context) error {
	return s.client.InvalidateCache(ctx, s.endpoint+"*")
}
