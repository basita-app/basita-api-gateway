package cms

import (
	"api-gateway/pkg/cache"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// CMSClient is the main client for interacting with Strapi CMS via GraphQL
type CMSClient struct {
	baseURL      string
	mediaBaseURL string // Base URL for media files (without /api)
	token        string
	httpClient   *http.Client
	cache        cache.Cache
	defaultTTL   time.Duration
}

// Config holds configuration for the CMS client
type Config struct {
	BaseURL         string
	MediaBaseURL    string // Optional: base URL for media (if different from BaseURL)
	Token           string
	RequestTimeout  time.Duration
	Cache           cache.Cache
	DefaultCacheTTL time.Duration
}

// NewCMSClient creates a new CMS client with the given configuration
func NewCMSClient(config Config) *CMSClient {
	if config.RequestTimeout == 0 {
		config.RequestTimeout = 30 * time.Second
	}
	if config.DefaultCacheTTL == 0 {
		config.DefaultCacheTTL = 5 * time.Minute
	}
	if config.Cache == nil {
		config.Cache = &cache.NoOpCache{}
	}

	// If MediaBaseURL is not provided, derive it from BaseURL by removing /api suffix
	mediaBaseURL := config.MediaBaseURL
	if mediaBaseURL == "" {
		mediaBaseURL = strings.TrimSuffix(config.BaseURL, "/graphql")
	}

	return &CMSClient{
		baseURL:      config.BaseURL,
		mediaBaseURL: mediaBaseURL,
		token:        config.Token,
		httpClient: &http.Client{
			Timeout: config.RequestTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		cache:      config.Cache,
		defaultTTL: config.DefaultCacheTTL,
	}
}

// InvalidateCache invalidates cache entries matching the given pattern
func (c *CMSClient) InvalidateCache(ctx context.Context, pattern string) error {
	return c.cache.DeletePattern(ctx, pattern)
}

// GraphQLRequest represents a GraphQL query request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message string        `json:"message"`
	Path    []interface{} `json:"path,omitempty"`
}

// ExecuteGraphQL executes a GraphQL query with caching support
func (c *CMSClient) ExecuteGraphQL(ctx context.Context, query string, variables map[string]interface{}) (json.RawMessage, error) {
	// Build cache key from query and variables
	cacheKey := c.buildCacheKey(query, variables)

	// Check cache first
	if cached, err := c.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		fmt.Printf("[GraphQL Client] Cache hit for key: %s\n", cacheKey)
		return cached, nil
	}
	// Create request body
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create GraphQL request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	// Debug logging
	fmt.Printf("[GraphQL Client] POST %s\n", c.baseURL)
	fmt.Printf("[GraphQL Client] Query: %s\n", query)
	if variables != nil {
		varsJSON, _ := json.Marshal(variables)
		fmt.Printf("[GraphQL Client] Variables: %s\n", string(varsJSON))
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GraphQL request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GraphQL response: %w", err)
	}

	// Debug logging
	fmt.Printf("[GraphQL Client] Response Status: %d\n", resp.StatusCode)
	fmt.Printf("[GraphQL Client] Response Body: %s\n", string(body))

	// Parse GraphQL response
	var graphqlResp GraphQLResponse
	if err := json.Unmarshal(body, &graphqlResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GraphQL response: %w", err)
	}

	// Check for GraphQL errors
	if len(graphqlResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", graphqlResp.Errors)
	}

	// Cache successful response
	if graphqlResp.Data != nil {
		_ = c.cache.Set(ctx, cacheKey, graphqlResp.Data, c.defaultTTL)
		fmt.Printf("[GraphQL Client] Cached response with key: %s\n", cacheKey)
	}

	return graphqlResp.Data, nil
}

// buildCacheKey generates a cache key from the query and variables
func (c *CMSClient) buildCacheKey(query string, variables map[string]interface{}) string {
	// Create a deterministic string representation
	var keyParts []string
	keyParts = append(keyParts, query)

	if variables != nil {
		// Marshal variables to JSON for consistent key generation
		varsJSON, _ := json.Marshal(variables)
		keyParts = append(keyParts, string(varsJSON))
	}

	// Combine and hash to create a shorter key
	combined := strings.Join(keyParts, "|")
	hash := sha256.Sum256([]byte(combined))
	return fmt.Sprintf("cms:graphql:%x", hash)
}

// PrefixMediaURL adds the media base URL prefix to a relative URL
// Returns empty string if url is empty, returns as-is if already absolute
func (c *CMSClient) PrefixMediaURL(url string) string {
	if url == "" {
		return ""
	}
	// If URL is already absolute, return as-is
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	// Add prefix to relative URL
	return c.mediaBaseURL + url
}

// strapiMediaField represents the raw media field structure from Strapi GraphQL
// This is used internally by services to parse media fields from GraphQL responses
type strapiMediaField struct {
	DocumentID string                 `json:"documentId"`
	URL        string                 `json:"url"`
	Width      int                    `json:"width"`
	Height     int                    `json:"height"`
	Formats    map[string]interface{} `json:"formats"`
}
