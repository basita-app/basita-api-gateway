package cms

import (
	"api-gateway/services/cms/models"
	"encoding/json"
	"fmt"
)

// CacheKeyBuilder helps build consistent cache keys for CMS resources
type CacheKeyBuilder struct {
	resource string
	parts    []string
}

// NewCacheKeyBuilder creates a new cache key builder
func NewCacheKeyBuilder(resource string) *CacheKeyBuilder {
	return &CacheKeyBuilder{
		resource: resource,
		parts:    make([]string, 0),
	}
}

// AddPart adds a part to the cache key
func (b *CacheKeyBuilder) AddPart(part string) *CacheKeyBuilder {
	if part != "" {
		b.parts = append(b.parts, part)
	}
	return b
}

// AddOptions adds collection options to the cache key
func (b *CacheKeyBuilder) AddOptions(opts models.CollectionOptions) *CacheKeyBuilder {
	if opts.Page > 0 {
		b.AddPart(fmt.Sprintf("page:%d", opts.Page))
	}
	if opts.PageSize > 0 {
		b.AddPart(fmt.Sprintf("pageSize:%d", opts.PageSize))
	}
	if opts.Populate != "" {
		b.AddPart(fmt.Sprintf("populate:%s", opts.Populate))
	}
	if opts.Locale != "" {
		b.AddPart(fmt.Sprintf("locale:%s", opts.Locale))
	}
	if len(opts.Sort) > 0 {
		sortJSON, _ := json.Marshal(opts.Sort)
		b.AddPart(fmt.Sprintf("sort:%s", sortJSON))
	}
	if len(opts.Filters) > 0 {
		filtersJSON, _ := json.Marshal(opts.Filters)
		b.AddPart(fmt.Sprintf("filters:%s", filtersJSON))
	}
	if len(opts.Fields) > 0 {
		fieldsJSON, _ := json.Marshal(opts.Fields)
		b.AddPart(fmt.Sprintf("fields:%s", fieldsJSON))
	}
	return b
}

// Build constructs the final cache key with cms_ prefix
func (b *CacheKeyBuilder) Build() string {
	key := "cms_" + b.resource
	for _, part := range b.parts {
		key += ":" + part
	}
	return key
}
