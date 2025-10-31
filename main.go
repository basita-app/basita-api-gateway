package main

import (
	"api-gateway/pkg/cache"
	"api-gateway/services/cms"
	"api-gateway/services/cms/models"
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "API Gateway",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "api-gateway",
		})
	})

	// Initialize Redis client for caching
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0 // Use DB 0 for CMS cache

	redisClient := redis.NewClient(&redis.Options{
		Addr:            redisAddr,
		Password:        redisPassword,
		DB:              redisDB,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v. Running without cache.", err)
		redisClient = nil // Disable cache if Redis is unavailable
	} else {
		log.Println("Redis connection established")
	}

	// Initialize CMS client with Redis cache
	cmsServiceURL := os.Getenv("CMS_SERVICE_URL")
	if cmsServiceURL == "" {
		cmsServiceURL = "http://localhost:1337/api"
	}
	cmsServiceToken := os.Getenv("CMS_SERVICE_TOKEN")

	// Debug: Check if token is loaded
	if cmsServiceToken == "" {
		log.Println("WARNING: CMS_SERVICE_TOKEN is not set!")
	}

	var cmsCache cache.Cache
	if redisClient != nil {
		cmsCache = cache.NewRedisCache(redisClient, "cms:")
	} else {
		cmsCache = &cache.NoOpCache{} // Fallback to no-op cache
	}

	cmsClient := cms.NewCMSClient(cms.Config{
		BaseURL:         cmsServiceURL,
		Token:           cmsServiceToken,
		RequestTimeout:  10 * time.Second,
		Cache:           cmsCache,
		DefaultCacheTTL: 24 * time.Hour,
	})

	// Initialize CMS services
	brandService := cms.NewBrandService(cmsClient)
	advertisementService := cms.NewAdvertisementService(cmsClient)
	_ = cms.NewCarModelService(cmsClient)
	_ = cms.NewCarVariantService(cmsClient)
	_ = cms.NewCityService(cmsClient)
	_ = cms.NewGovernorateService(cmsClient)
	_ = cms.NewShowroomService(cmsClient)

	app.Get("/uploads/:path", func(c *fiber.Ctx) error {
		path := c.Params("path")
		cmsUrlWithoutAPI := strings.TrimSuffix(cmsServiceURL, "/api")
		println(cmsUrlWithoutAPI)
		url := cmsUrlWithoutAPI + "/uploads/" + path
		return proxy.Do(c, url)
	})
	// CMS API Routes
	cmsGroup := app.Group("/api/cms")

	// Brands endpoints
	cmsGroup.Get("/brands", func(c *fiber.Ctx) error {
		ctx := c.Context()

		// Extract locale from Accept-Language header
		// Keep the full locale code (e.g., "ar-EG", "en")
		locale := c.Get("Accept-Language")
		// If multiple locales are provided (e.g., "en-US,en;q=0.9"), take the first one
		if locale != "" {
			if commaIdx := strings.Index(locale, ","); commaIdx != -1 {
				locale = locale[:commaIdx]
			}
			locale = strings.TrimSpace(locale)
		}

		opts := models.CollectionOptions{
			Populate: c.Query("populate", "*"),
			Locale:   locale,
		}

		brands, err := brandService.GetSimplified(ctx, opts, true)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(brands)
	})

	cmsGroup.Get("/brands/:id", func(c *fiber.Ctx) error {
		ctx := c.Context()
		id := c.Params("id")

		// Extract locale from Accept-Language header
		// Keep the full locale code (e.g., "ar-EG", "en")
		locale := c.Get("Accept-Language")
		// If multiple locales are provided (e.g., "en-US,en;q=0.9"), take the first one
		if locale != "" {
			if commaIdx := strings.Index(locale, ","); commaIdx != -1 {
				locale = locale[:commaIdx]
			}
			locale = strings.TrimSpace(locale)
		}

		opts := models.ItemOptions{
			Populate: c.Query("populate", "*"),
			Locale:   locale,
		}

		brand, err := brandService.GetByID(ctx, id, opts, true)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Brand not found",
			})
		}

		return c.JSON(brand)
	})

	// Advertisements endpoints
	cmsGroup.Get("/advertisements", func(c *fiber.Ctx) error {
		ctx := c.Context()

		// Extract locale from Accept-Language header
		// Keep the full locale code (e.g., "ar-EG", "en")
		locale := c.Get("Accept-Language")
		// If multiple locales are provided (e.g., "en-US,en;q=0.9"), take the first one
		if locale != "" {
			if commaIdx := strings.Index(locale, ","); commaIdx != -1 {
				locale = locale[:commaIdx]
			}
			locale = strings.TrimSpace(locale)
		}

		opts := models.CollectionOptions{
			Page:     c.QueryInt("page", 1),
			PageSize: c.QueryInt("pageSize", 25),
			Populate: c.Query("populate", "*"),
			Locale:   locale,
		}

		advertisements, err := advertisementService.GetAll(ctx, opts, true)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(advertisements)
	})

	cmsGroup.Get("/advertisements/:id", func(c *fiber.Ctx) error {
		ctx := c.Context()
		id := c.Params("id")

		// Extract locale from Accept-Language header
		// Keep the full locale code (e.g., "ar-EG", "en")
		locale := c.Get("Accept-Language")
		// If multiple locales are provided (e.g., "en-US,en;q=0.9"), take the first one
		if locale != "" {
			if commaIdx := strings.Index(locale, ","); commaIdx != -1 {
				locale = locale[:commaIdx]
			}
			locale = strings.TrimSpace(locale)
		}

		opts := models.ItemOptions{
			Populate: c.Query("populate", "*"),
			Locale:   locale,
		}

		advertisement, err := advertisementService.GetByID(ctx, id, opts, true)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Advertisement not found",
			})
		}

		return c.JSON(advertisement)
	})

	// Cache Management Endpoint
	cacheSecretKey := os.Getenv("CACHE_SECRET_KEY")
	app.Post("/api/cache/invalidate/cms", func(c *fiber.Ctx) error {
		// Check secret key
		secretKey := c.Get("X-Cache-Secret-Key")
		if secretKey == "" {
			secretKey = c.Query("secret")
		}

		if cacheSecretKey == "" || secretKey != cacheSecretKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid or missing secret key",
			})
		}

		// Get pattern from query or body
		pattern := c.Query("pattern")
		if pattern == "" {
			pattern = "cms_*" // Default: clear all CMS cache
		}

		// Invalidate cache
		ctx := c.Context()
		if err := cmsClient.InvalidateCache(ctx, pattern); err != nil {
			log.Printf("Failed to invalidate cache: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to invalidate cache",
				"details": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Cache invalidated successfully",
			"pattern": pattern,
		})
	})

	// Proxy routes
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:3000"
	}

	app.All("/api/auth/*", func(c *fiber.Ctx) error {
		url := authServiceURL + c.Path()
		if len(c.Request().URI().QueryString()) > 0 {
			url += "?" + string(c.Request().URI().QueryString())
		}

		if err := proxy.Do(c, url); err != nil {
			return err
		}

		c.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	carListingServiceURL := os.Getenv("CAR_LISTING_SERVICE_URL")
	if carListingServiceURL == "" {
		carListingServiceURL = "http://localhost:3002"
	}

	app.All("/api/car-listing/*", func(c *fiber.Ctx) error {
		url := carListingServiceURL + c.Path()
		if len(c.Request().URI().QueryString()) > 0 {
			url += "?" + string(c.Request().URI().QueryString())
		}

		if err := proxy.Do(c, url); err != nil {
			return err
		}

		c.Response().Header.Del(fiber.HeaderServer)
		return nil
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	// Start server
	log.Printf("API Gateway starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
