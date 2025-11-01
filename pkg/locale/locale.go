package locale

import (
	"context"
	"strings"
)

// contextKey is a private type for context keys to avoid collisions
type contextKey string

const localeContextKey contextKey = "locale"

// DefaultLocale is the fallback locale if none is specified
const DefaultLocale = "en"

// FromContext extracts the locale from the context
func FromContext(ctx context.Context) string {
	if locale, ok := ctx.Value(localeContextKey).(string); ok && locale != "" {
		return locale
	}
	return DefaultLocale
}

// WithLocale adds a locale to the context
func WithLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, localeContextKey, locale)
}

// ParseAcceptLanguage parses the Accept-Language header and returns the primary locale
func ParseAcceptLanguage(acceptLanguage string) string {
	if acceptLanguage == "" {
		return DefaultLocale
	}

	// Accept-Language format: "en-US,en;q=0.9,ar;q=0.8"
	// We take the first locale (highest priority)
	parts := strings.Split(acceptLanguage, ",")
	if len(parts) == 0 {
		return DefaultLocale
	}

	// Get the first locale and remove quality value if present
	firstLocale := strings.TrimSpace(parts[0])
	firstLocale = strings.Split(firstLocale, ";")[0]

	return firstLocale
}
