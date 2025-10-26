package main

import (
	"api-gateway/pkg/cache"
	"api-gateway/services/cms"
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
)

func main() {
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
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
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

	var cmsCache cache.Cache
	if redisClient != nil {
		cmsCache = cache.NewRedisCache(redisClient, "cms:")
	} else {
		cmsCache = &cache.NoOpCache{} // Fallback to no-op cache
	}

	cmsClient := cms.NewCMSClient(cms.Config{
		BaseURL:         cmsServiceURL,
		Token:           cmsServiceToken,
		RequestTimeout:  30 * time.Second,
		Cache:           cmsCache,
		DefaultCacheTTL: 10 * time.Minute,
	})

	// Initialize CMS services (example - add your actual services here)
	_ = cmsClient // Use cmsClient to initialize your specific resource services
	// Example:
	// articleService := cms.NewArticleService(cmsClient)
	// homepageService := cms.NewHomepageService(cmsClient)

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
