package main

import (
	"api-gateway/pkg/cache"
	"api-gateway/pkg/locale"
	"api-gateway/services/cms"
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
	"github.com/goccy/go-json"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "API Gateway",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Locale middleware - extracts Accept-Language header and adds to context
	app.Use(func(c *fiber.Ctx) error {
		acceptLanguage := c.Get("Accept-Language", "en")
		userLocale := locale.ParseAcceptLanguage(acceptLanguage)
		log.Printf("locale %s", userLocale)

		// Add locale to the user context
		ctx := c.UserContext()
		ctx = locale.WithLocale(ctx, userLocale)
		c.SetUserContext(ctx)

		return c.Next()
	})

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "api-gateway",
		})
	})

	// API Documentation
	app.Get("/api-docs", func(c *fiber.Ctx) error {
		return c.SendFile("./index.html")
	})

	app.Get("/openapi.yaml", func(c *fiber.Ctx) error {
		return c.SendFile("./openapi.yaml")
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
		cmsServiceURL = "http://localhost:1337/graphql"
	}
	cmsServiceToken := os.Getenv("CMS_SERVICE_TOKEN")

	cmsMediaBaseURL := os.Getenv("CMS_MEDIA_BASE_URL")
	if cmsMediaBaseURL == "" {
		cmsMediaBaseURL = "http://localhost:1337"
	}

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
		MediaBaseURL:   cmsMediaBaseURL,
		Token:           cmsServiceToken,
		RequestTimeout:  10 * time.Second,
		Cache:           cmsCache,
		DefaultCacheTTL: 24 * time.Hour,
	})

	// Initialize CMS GraphQL services
	brandService := cms.NewBrandServiceGraphQL(cmsClient)
	advertisementService := cms.NewAdvertisementServiceGraphQL(cmsClient)
	carModelService := cms.NewCarModelServiceGraphQL(cmsClient)
	showroomService := cms.NewShowroomServiceGraphQL(cmsClient)
	governorateService := cms.NewGovernorateServiceGraphQL(cmsClient)
	appVersionService := cms.NewAppVersionServiceGraphQL(cmsClient)
	_ = governorateService // Service initialized but endpoint not yet added
	_ = appVersionService  // Service initialized but endpoint not yet added

	app.Get("/uploads/:path", func(c *fiber.Ctx) error {
		path := c.Params("path")
		cmsUrlWithoutAPI := strings.TrimSuffix(cmsServiceURL, "/graphql")
		println(cmsUrlWithoutAPI)
		url := cmsUrlWithoutAPI + "/uploads/" + path
		return proxy.Do(c, url)
	})
	// CMS API Routes
	cmsGroup := app.Group("/api/cms")

	// Brands endpoints
	cmsGroup.Get("/brands", func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		brands, err := brandService.GetSimplified(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(brands)
	})

	cmsGroup.Get("/brands/:id", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		id := c.Params("id")

		brand, err := brandService.GetByID(ctx, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Brand not found",
			})
		}

		return c.JSON(brand)
	})

	// Get car models by brand ID
	cmsGroup.Get("/brands/:id/cars", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		brandID := c.Params("id")

		carModels, err := carModelService.GetByBrandID(ctx, brandID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(carModels)
	})

	// Get detailed car model by ID
	cmsGroup.Get("/cars/:id", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		id := c.Params("id")

		carModel, err := carModelService.GetDetailedByID(ctx, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Car model not found",
			})
		}

		return c.JSON(carModel)
	})

	// Get detailed variant by ID
	cmsGroup.Get("/cars/:carId/variants/:variantId", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		variantID := c.Params("variantId")

		variant, err := carModelService.GetVariantByID(ctx, variantID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Variant not found",
			})
		}

		return c.JSON(variant)
	})

	// Advertisements endpoints
	cmsGroup.Get("/advertisements", func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		advertisements, err := advertisementService.GetAll(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(advertisements)
	})

	cmsGroup.Get("/advertisements/:id", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		id := c.Params("id")

		advertisement, err := advertisementService.GetByID(ctx, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Advertisement not found",
			})
		}

		return c.JSON(advertisement)
	})

	// Showrooms endpoints
	cmsGroup.Get("/showrooms", func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		showrooms, err := showroomService.GetAll(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(showrooms)
	})

	cmsGroup.Get("/showrooms/:id", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		id := c.Params("id")

		showroom, err := showroomService.GetByID(ctx, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Showroom not found",
			})
		}

		return c.JSON(showroom)
	})

	// Get car variants by showroom ID
	cmsGroup.Get("/showrooms/:id/variants", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		showroomID := c.Params("id")

		variants, err := showroomService.GetCarVariantsByShowroomID(ctx, showroomID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(variants)
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
			pattern = "cms:graphql:*" // Default: clear all CMS GraphQL cache
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
