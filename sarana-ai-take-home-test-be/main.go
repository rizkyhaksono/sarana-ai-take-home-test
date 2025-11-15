package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/database"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/docs"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/middleware"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/routes"

	_ "github.com/rizkyhaksono/sarana-ai-take-home-test/docs"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Notes API
// @version 2.0
// @description A notes management API with user authentication and image upload capabilities
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@sarana-ai.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// if err := utils.InitLoki(); err != nil {
	// 	log.Println("Warning: Failed to initialize Loki client:", err)
	// 	log.Println("Continuing without Loki logging...")
	// }
	// defer utils.StopLoki()

	host := getEnv("SWAGGER_HOST", "40.90.171.103:36322")
	docs.SwaggerInfo.Host = host

	schemes := getEnv("SWAGGER_SCHEMES", "http")
	if schemes == "https" {
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	}

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

	app.Use(cors.New(cors.Config{
		AllowOrigins:     getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:8080"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-Request-ID, X-Client-Version",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length, Content-Type",
	}))

	app.Use(middleware.Logger())

	routes.SetupRoutes(app)

	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Get("/docs", func(c *fiber.Ctx) error {
		html := `<!doctype html>
<html>
<head>
    <title>API Documentation</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
    <script id="api-reference" data-url="/swagger/doc.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

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
