package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/handlers"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/middleware"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Public routes
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	// Protected routes - Notes
	api := app.Group("/notes", middleware.JWTAuth)

	api.Post("/", handlers.CreateNote)
	api.Get("/", handlers.GetNotes)
	api.Get("/:id", handlers.GetNote)
	api.Delete("/:id", handlers.DeleteNote)
	api.Post("/:id/image", handlers.UploadNoteImage)
	api.Get("/:id/image", handlers.GetNoteImage)

	// Protected routes - Logs
	logs := app.Group("/logs", middleware.JWTAuth)

	logs.Get("/", handlers.GetLogs)
	logs.Get("/:id", handlers.GetLog)

	app.Static("/uploads", "./uploads")
}
