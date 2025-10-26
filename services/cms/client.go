package cms

import (
	"api-gateway/pkg/cache"
	"api-gateway/services/cms/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the interface for CMS operations
type Client interface {
	GetCollection(ctx context.Context, endpoint string, opts models.CollectionOptions, cacheOpts *models.CacheOptions) ([]byte, error)
	GetItem(ctx context.Context, endpoint string, id string, opts models.ItemOptions, cacheOpts *models.CacheOptions) ([]byte, error)
	GetSingle(ctx context.Context, endpoint string, opts models.ItemOptions, cacheOpts *models.CacheOptions) ([]byte, error)
	InvalidateCache(ctx context.Context, pattern string) error
}

// CMSClient is the main client for interacting with Strapi CMS
type CMSClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
	cache      cache.Cache
	defaultTTL time.Duration
}

// Config holds configuration for the CMS client
type Config struct {
	BaseURL         string
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

	// Ensure baseURL ends with /
	baseURL := strings.TrimRight(config.BaseURL, "/") + "/"

	return &CMSClient{
		baseURL: baseURL,
		token:   config.Token,
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

// GetCollection fetches a collection from Strapi with caching support
func (c *CMSClient) GetCollection(ctx context.Context, endpoint string, opts models.CollectionOptions, cacheOpts *models.CacheOptions) ([]byte, error) {
	// Build cache key
	cacheKey := NewCacheKeyBuilder(endpoint).
		AddOptions(opts).
		Build()

	// Check cache if enabled
	if cacheOpts != nil && cacheOpts.Enabled {
		if cached, err := c.cache.Get(ctx, cacheKey); err == nil && cached != nil {
			return cached, nil
		}
	}

	// Build request URL
	reqURL := c.buildCollectionURL(endpoint, opts)

	// Make request
	data, err := c.doRequest(ctx, "GET", reqURL)
	if err != nil {
		return nil, err
	}

	// Cache response if enabled
	if cacheOpts != nil && cacheOpts.Enabled {
		ttl := c.defaultTTL
		if cacheOpts.TTL > 0 {
			ttl = cacheOpts.TTL
		}
		// We don't return error if caching fails, just log it
		_ = c.cache.Set(ctx, cacheKey, data, ttl)
	}

	return data, nil
}

// GetItem fetches a single item by ID from Strapi with caching support
func (c *CMSClient) GetItem(ctx context.Context, endpoint string, id string, opts models.ItemOptions, cacheOpts *models.CacheOptions) ([]byte, error) {
	// Build cache key
	cacheKey := NewCacheKeyBuilder(endpoint).
		AddPart(id).
		AddPart(fmt.Sprintf("populate:%s", opts.Populate)).
		AddPart(fmt.Sprintf("locale:%s", opts.Locale)).
		Build()

	// Check cache if enabled
	if cacheOpts != nil && cacheOpts.Enabled {
		if cached, err := c.cache.Get(ctx, cacheKey); err == nil && cached != nil {
			return cached, nil
		}
	}

	// Build request URL
	reqURL := c.buildItemURL(endpoint, id, opts)

	// Make request
	data, err := c.doRequest(ctx, "GET", reqURL)
	if err != nil {
		return nil, err
	}

	// Cache response if enabled
	if cacheOpts != nil && cacheOpts.Enabled {
		ttl := c.defaultTTL
		if cacheOpts.TTL > 0 {
			ttl = cacheOpts.TTL
		}
		_ = c.cache.Set(ctx, cacheKey, data, ttl)
	}

	return data, nil
}

// GetSingle fetches a single-type resource from Strapi with caching support
func (c *CMSClient) GetSingle(ctx context.Context, endpoint string, opts models.ItemOptions, cacheOpts *models.CacheOptions) ([]byte, error) {
	// Build cache key
	cacheKey := NewCacheKeyBuilder(endpoint).
		AddPart(fmt.Sprintf("populate:%s", opts.Populate)).
		AddPart(fmt.Sprintf("locale:%s", opts.Locale)).
		Build()

	// Check cache if enabled
	if cacheOpts != nil && cacheOpts.Enabled {
		if cached, err := c.cache.Get(ctx, cacheKey); err == nil && cached != nil {
			return cached, nil
		}
	}

	// Build request URL
	reqURL := c.buildSingleURL(endpoint, opts)

	// Make request
	data, err := c.doRequest(ctx, "GET", reqURL)
	if err != nil {
		return nil, err
	}

	// Cache response if enabled
	if cacheOpts != nil && cacheOpts.Enabled {
		ttl := c.defaultTTL
		if cacheOpts.TTL > 0 {
			ttl = cacheOpts.TTL
		}
		_ = c.cache.Set(ctx, cacheKey, data, ttl)
	}

	return data, nil
}

// InvalidateCache invalidates cache entries matching the given pattern
func (c *CMSClient) InvalidateCache(ctx context.Context, pattern string) error {
	return c.cache.DeletePattern(ctx, pattern)
}

// doRequest performs an HTTP request with proper headers and error handling
func (c *CMSClient) doRequest(ctx context.Context, method, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, &models.RequestError{
			Message: fmt.Sprintf("failed to create request: %v", err),
		}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &models.RequestError{
			Message: fmt.Sprintf("request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &models.RequestError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("failed to read response body: %v", err),
		}
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		var strapiErr models.StrapiError
		if err := json.Unmarshal(body, &strapiErr); err == nil && strapiErr.Error.Message != "" {
			return nil, &models.RequestError{
				StatusCode: resp.StatusCode,
				StrapiErr:  &strapiErr,
			}
		}
		return nil, &models.RequestError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("request failed with status %d: %s", resp.StatusCode, string(body)),
		}
	}

	return body, nil
}

// buildCollectionURL constructs the URL for collection requests
func (c *CMSClient) buildCollectionURL(endpoint string, opts models.CollectionOptions) string {
	u, _ := url.Parse(c.baseURL + strings.TrimPrefix(endpoint, "/"))
	query := u.Query()

	if opts.Page > 0 {
		query.Set("pagination[page]", fmt.Sprintf("%d", opts.Page))
	}
	if opts.PageSize > 0 {
		query.Set("pagination[pageSize]", fmt.Sprintf("%d", opts.PageSize))
	}
	if opts.Populate != "" {
		query.Set("populate", opts.Populate)
	}
	if opts.Locale != "" {
		query.Set("locale", opts.Locale)
	}

	// Add filters
	for key, value := range opts.Filters {
		query.Set(key, value)
	}

	// Add sort
	for _, sort := range opts.Sort {
		query.Add("sort", sort)
	}

	// Add fields
	for _, field := range opts.Fields {
		query.Add("fields", field)
	}

	u.RawQuery = query.Encode()
	return u.String()
}

// buildItemURL constructs the URL for single item requests
func (c *CMSClient) buildItemURL(endpoint, id string, opts models.ItemOptions) string {
	u, _ := url.Parse(c.baseURL + strings.TrimPrefix(endpoint, "/") + "/" + id)
	query := u.Query()

	if opts.Populate != "" {
		query.Set("populate", opts.Populate)
	}
	if opts.Locale != "" {
		query.Set("locale", opts.Locale)
	}

	// Add fields
	for _, field := range opts.Fields {
		query.Add("fields", field)
	}

	u.RawQuery = query.Encode()
	return u.String()
}

// buildSingleURL constructs the URL for single-type requests
func (c *CMSClient) buildSingleURL(endpoint string, opts models.ItemOptions) string {
	u, _ := url.Parse(c.baseURL + strings.TrimPrefix(endpoint, "/"))
	query := u.Query()

	if opts.Populate != "" {
		query.Set("populate", opts.Populate)
	}
	if opts.Locale != "" {
		query.Set("locale", opts.Locale)
	}

	// Add fields
	for _, field := range opts.Fields {
		query.Add("fields", field)
	}

	u.RawQuery = query.Encode()
	return u.String()
}
