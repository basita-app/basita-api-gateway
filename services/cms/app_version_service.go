package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"fmt"
)

// AppVersionServiceGraphQL handles operations for application version using GraphQL
type AppVersionServiceGraphQL struct {
	client *CMSClient
}

// NewAppVersionServiceGraphQL creates a new GraphQL-based service for application version
func NewAppVersionServiceGraphQL(client *CMSClient) *AppVersionServiceGraphQL {
	return &AppVersionServiceGraphQL{
		client: client,
	}
}

// Get fetches the current application version information
func (s *AppVersionServiceGraphQL) Get(ctx context.Context) (*models.ApplicationVersion, error) {
	data, err := s.client.ExecuteGraphQL(ctx, GetAppVersionQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch application version: %w", err)
	}

	var result struct {
		ApplicationVersion *struct {
			MobileAppVersion     string `json:"MobileAppVersion"`
			MobileAppBuildNumber string `json:"MobileAppBuildNumber"`
			WebVersion           string `json:"WebVersion"`
		} `json:"applicationVersion"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal application version: %w", err)
	}

	if result.ApplicationVersion == nil {
		return nil, fmt.Errorf("application version not found")
	}

	return &models.ApplicationVersion{
		MobileAppVersion:     result.ApplicationVersion.MobileAppVersion,
		MobileAppBuildNumber: result.ApplicationVersion.MobileAppBuildNumber,
		WebVersion:           result.ApplicationVersion.WebVersion,
	}, nil
}
