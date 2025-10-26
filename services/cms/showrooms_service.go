package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"time"
)

// ShowroomService handles operations for the showrooms resource
type ShowroomService struct {
	client   Client
	endpoint string
	cacheTTL time.Duration
}

// NewShowroomService creates a new service for showrooms
func NewShowroomService(client Client) *ShowroomService {
	return &ShowroomService{
		client:   client,
		endpoint: "showrooms",
		cacheTTL: 10 * time.Minute,
	}
}

// GetAll fetches all showrooms with optional filters
func (s *ShowroomService) GetAll(ctx context.Context, opts models.CollectionOptions, useCache bool) (*models.ShowroomCollectionResponse, error) {
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

	var response models.ShowroomCollectionResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal showrooms: " + err.Error(),
		}
	}

	return &response, nil
}

// GetByID fetches a single showroom by ID
func (s *ShowroomService) GetByID(ctx context.Context, id string, opts models.ItemOptions, useCache bool) (*models.ShowroomResponse, error) {
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

	var response models.ShowroomResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, &models.RequestError{
			Message: "failed to unmarshal showroom: " + err.Error(),
		}
	}

	return &response, nil
}

// InvalidateCache invalidates all cached showrooms
func (s *ShowroomService) InvalidateCache(ctx context.Context) error {
	return s.client.InvalidateCache(ctx, s.endpoint+"*")
}
