package cms

import (
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"fmt"
)

// GovernorateServiceGraphQL handles operations for governorates using GraphQL
type GovernorateServiceGraphQL struct {
	client *CMSClient
}

// NewGovernorateServiceGraphQL creates a new GraphQL-based service for governorates
func NewGovernorateServiceGraphQL(client *CMSClient) *GovernorateServiceGraphQL {
	return &GovernorateServiceGraphQL{
		client: client,
	}
}

// GetAll fetches all governorates with their cities
func (s *GovernorateServiceGraphQL) GetAll(ctx context.Context) ([]models.GovernorateWithCities, error) {
	data, err := s.client.ExecuteGraphQL(ctx, GetGovernoratesQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch governorates: %w", err)
	}

	var result struct {
		Governorates []struct {
			DocumentID string `json:"documentId"`
			Name       string `json:"Name"`
			Cities     []struct {
				DocumentID string `json:"documentId"`
				Name       string `json:"Name"`
			} `json:"cities"`
		} `json:"governorates"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal governorates: %w", err)
	}

	governorates := make([]models.GovernorateWithCities, len(result.Governorates))
	for i, gov := range result.Governorates {
		cities := make([]models.City, len(gov.Cities))
		for j, city := range gov.Cities {
			cities[j] = models.City{
				ID:   city.DocumentID,
				Name: city.Name,
			}
		}

		governorates[i] = models.GovernorateWithCities{
			ID:     gov.DocumentID,
			Name:   gov.Name,
			Cities: cities,
		}
	}

	return governorates, nil
}
