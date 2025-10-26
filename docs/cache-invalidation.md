# Cache Invalidation API

## Endpoint

```
POST /api/cache/invalidate/cms
```

Protected endpoint for invalidating CMS cache entries.

## Authentication

Requires a secret key set in the `CACHE_SECRET_KEY` environment variable.

The secret key can be provided in two ways:

1. **Header**: `X-Cache-Secret-Key: your-secret-key`
2. **Query parameter**: `?secret=your-secret-key`

## Parameters

| Parameter | Type   | Required | Default | Description |
|-----------|--------|----------|---------|-------------|
| pattern   | string | No       | `cms_*` | Redis glob pattern for keys to invalidate |
| secret    | string | Yes      | -       | Secret key for authentication (if not in header) |

## Examples

### Clear All CMS Cache

```bash
curl -X POST "http://localhost:3001/api/cache/invalidate/cms" \
  -H "X-Cache-Secret-Key: your-secret-key"
```

### Clear Specific Resource

```bash
# Clear all articles cache
curl -X POST "http://localhost:3001/api/cache/invalidate/cms?pattern=cms_articles:*" \
  -H "X-Cache-Secret-Key: your-secret-key"
```

### Clear Specific Item

```bash
# Clear a specific article
curl -X POST "http://localhost:3001/api/cache/invalidate/cms?pattern=cms_articles:123:*" \
  -H "X-Cache-Secret-Key: your-secret-key"
```

### Using Query Parameter for Secret

```bash
curl -X POST "http://localhost:3001/api/cache/invalidate/cms?secret=your-secret-key&pattern=cms_products:*"
```

## Response

### Success (200)

```json
{
  "success": true,
  "message": "Cache invalidated successfully",
  "pattern": "cms_*"
}
```

### Unauthorized (401)

```json
{
  "error": "Unauthorized: Invalid or missing secret key"
}
```

### Error (500)

```json
{
  "error": "Failed to invalidate cache",
  "details": "error details here"
}
```

## Cache Key Patterns

All CMS cache keys are prefixed with `cms_` followed by the resource name.

### Pattern Examples

| Pattern | Description |
|---------|-------------|
| `cms_*` | All CMS cache entries (default) |
| `cms_articles:*` | All articles cache |
| `cms_articles:123:*` | Specific article with ID 123 |
| `cms_products:*` | All products cache |
| `cms_homepage:*` | Homepage cache |

## Common Use Cases

### 1. Strapi Webhook Integration

Configure Strapi to call this endpoint when content is published/updated:

**Strapi Webhook Configuration:**
- URL: `https://your-api-gateway.com/api/cache/invalidate/cms?secret=your-secret-key&pattern=cms_articles:*`
- Events: `entry.create`, `entry.update`, `entry.delete`, `entry.publish`, `entry.unpublish`

### 2. Manual Cache Clear

When you need to manually clear cache after bulk updates:

```bash
# Clear everything
curl -X POST "https://your-api-gateway.com/api/cache/invalidate/cms" \
  -H "X-Cache-Secret-Key: your-secret-key"
```

### 3. CI/CD Integration

Clear cache after deployment:

```bash
# In your deployment script
curl -X POST "$API_GATEWAY_URL/api/cache/invalidate/cms" \
  -H "X-Cache-Secret-Key: $CACHE_SECRET_KEY"
```

## Security Best Practices

1. **Use Strong Secret Key**: Generate a long, random secret key
   ```bash
   openssl rand -hex 32
   ```

2. **Use HTTPS**: Always use HTTPS in production

3. **Environment Variables**: Never commit the secret key to version control

4. **Header vs Query**: Prefer using the header method as query parameters may be logged

5. **Rotate Keys**: Periodically rotate your secret key

## Configuration

Add to your `.env` file:

```env
CACHE_SECRET_KEY=your-strong-random-secret-key-here
```

Generate a secure key:
```bash
# Generate a 32-byte random key
openssl rand -hex 32

# Or use Node.js
node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"
```

## Redis Pattern Matching

The `pattern` parameter uses Redis glob-style patterns:

- `*` - Matches any characters
- `?` - Matches exactly one character
- `[abc]` - Matches a, b, or c
- `[a-z]` - Matches any character from a to z

### Examples

```bash
# All articles
cms_articles:*

# Articles with specific page
cms_articles:page:1:*

# Articles in English locale
cms_articles:*:locale:en:*

# Everything starting with "cms_prod"
cms_prod*
```

## Error Handling

The endpoint will return an error if:

1. Secret key is missing or invalid (401)
2. Redis is unavailable (500)
3. Pattern is invalid (500)

Always check the response status code and handle errors appropriately.

## Monitoring

Consider logging all cache invalidation requests for audit purposes. The endpoint already logs failures to the server console.

## Rate Limiting

Consider implementing rate limiting for this endpoint to prevent abuse:

```go
// Example using fiber rate limiter
app.Use("/api/cache/invalidate/*", limiter.New(limiter.Config{
    Max:        10,
    Expiration: 1 * time.Minute,
}))
```
