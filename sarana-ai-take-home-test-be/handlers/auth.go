package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/constants"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/models"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/services"
)

var authService = services.NewAuthService()

// Register handles user registration
func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidRequestBody,
		})
	}

	// Register user
	user, token, err := authService.Register(req.Email, req.Password)
	if err != nil {
		switch err.Error() {
		case constants.ErrEmailExists:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		case constants.ErrHashingPassword, constants.ErrGeneratingToken:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": constants.ErrCreatingUser,
			})
		}
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Login handles user authentication
func Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidRequestBody,
		})
	}

	// Authenticate user
	user, token, err := authService.Login(req.Email, req.Password)
	if err != nil {
		switch err.Error() {
		case constants.ErrInvalidCredentials:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": constants.ErrUserNotFound,
			})
		}
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(models.AuthResponse{
		Token: token,
		User:  *user,
	})
}
