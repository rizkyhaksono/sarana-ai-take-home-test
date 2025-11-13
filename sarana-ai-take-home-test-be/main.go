package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/database"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/handlers"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/middleware"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/utils"
)

func main() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize database schema
	if err := database.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Initialize Loki client
	if err := utils.InitLoki(); err != nil {
		log.Println("Warning: Failed to initialize Loki client:", err)
		log.Println("Continuing without Loki logging...")
	}
	defer utils.StopLoki()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Apply logging middleware globally
	app.Use(middleware.Logger())

	// Serve static files (uploaded images)
	app.Static("/uploads", "./uploads")

	// Public routes
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	// Protected routes
	api := app.Group("/", middleware.JWTAuth)

	// Notes endpoints
	api.Post("/notes", handlers.CreateNote)
	api.Get("/notes", handlers.GetNotes)
	api.Get("/notes/:id", handlers.GetNote)
	api.Delete("/notes/:id", handlers.DeleteNote)
	api.Post("/notes/:id/image", handlers.UploadNoteImage)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
