package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
