package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/constants"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/models"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/services"
)

var authService = services.NewAuthService()

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration credentials"
// @Success 201 {object} models.BaseResponse "User registered successfully"
// @Failure 400 {object} models.BaseResponse "Invalid request body"
// @Failure 409 {object} models.BaseResponse "Email already exists"
// @Failure 500 {object} models.BaseResponse "Internal server error"
// @Router /register [post]
func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.ErrorResponse("INVALID_REQUEST", constants.ErrInvalidRequestBody, err.Error()),
		)
	}

	user, token, err := authService.Register(req.Email, req.Password)
	if err != nil {
		switch err.Error() {
		case constants.ErrEmailExists:
			return c.Status(fiber.StatusConflict).JSON(
				models.ErrorResponse("EMAIL_EXISTS", err.Error(), "The email address is already registered"),
			)
		case constants.ErrHashingPassword:
			return c.Status(fiber.StatusInternalServerError).JSON(
				models.ErrorResponse("HASHING_ERROR", err.Error(), "Failed to hash password"),
			)
		case constants.ErrGeneratingToken:
			return c.Status(fiber.StatusInternalServerError).JSON(
				models.ErrorResponse("TOKEN_ERROR", err.Error(), "Failed to generate authentication token"),
			)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(
				models.ErrorResponse("CREATE_USER_ERROR", constants.ErrCreatingUser, err.Error()),
			)
		}
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(
		models.SuccessResponse("User registered successfully", fiber.Map{
			"token": token,
			"user":  user,
		}),
	)
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.BaseResponse "Login successful"
// @Failure 400 {object} models.BaseResponse "Invalid request body"
// @Failure 401 {object} models.BaseResponse "Invalid credentials"
// @Failure 500 {object} models.BaseResponse "Internal server error"
// @Router /login [post]
func Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.ErrorResponse("INVALID_REQUEST", constants.ErrInvalidRequestBody, err.Error()),
		)
	}

	user, token, err := authService.Login(req.Email, req.Password)
	if err != nil {
		switch err.Error() {
		case constants.ErrInvalidCredentials:
			return c.Status(fiber.StatusUnauthorized).JSON(
				models.ErrorResponse("INVALID_CREDENTIALS", err.Error(), "Email or password is incorrect"),
			)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(
				models.ErrorResponse("USER_NOT_FOUND", constants.ErrUserNotFound, err.Error()),
			)
		}
	}

	return c.Status(fiber.StatusOK).JSON(
		models.SuccessResponse("Login successful", fiber.Map{
			"token": token,
			"user":  user,
		}),
	)
}

// Profile handles user authentication
// @Summary Profile user
// @Description Authenticate user with bearertoken
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} models.BaseResponse "User profile fetched successfully"
// @Failure 400 {object} models.BaseResponse "Invalid request body"
// @Failure 401 {object} models.BaseResponse "Unauthorized"
// @Failure 500 {object} models.BaseResponse "Internal server error"
// @Router /me [get]
func GetUserProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(
			models.ErrorResponse("INVALID_TOKEN", constants.ErrInvalidToken, "User ID not found in context"),
		)
	}

	user, err := authService.UserProfile(userID)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			models.ErrorResponse("INVALID_TOKEN", constants.ErrInvalidToken, err.Error()),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		models.SuccessResponse("User profile fetched successfully", fiber.Map{
			"user": user,
		}),
	)
}
