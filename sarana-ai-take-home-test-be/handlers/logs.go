package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/models"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/services"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/utils"
)

var logService = services.NewLogService()

// GetLogs retrieves all logs with search, sort, and pagination
// @Summary Get all logs
// @Description Retrieve all application logs with optional search, sorting, and pagination
// @Tags Logs
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search in method, endpoint, request_body, response_body"
// @Param sort_by query string false "Sort by field (datetime, created_at, method, endpoint, status_code)" default(datetime)
// @Param order query string false "Sort order (ASC, DESC)" default(DESC)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} services.LogsResponse "Paginated list of logs"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /logs [get]
func GetLogs(c *fiber.Ctx) error {
	params := utils.PaginationParams{
		Search: c.Query("search", ""),
		SortBy: c.Query("sort_by", "datetime"),
		Order:  c.Query("order", "DESC"),
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
	}

	result, err := logService.GetLogsWithParams(params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.ErrorResponse("GET_LOGS_ERROR", "Failed to retrieve logs", err.Error()),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		models.SuccessResponse("Logs retrieved successfully", result),
	)
}

// GetLog retrieves a single log by ID
// @Summary Get a log by ID
// @Description Retrieve a specific log entry by its ID
// @Tags Logs
// @Produce json
// @Security BearerAuth
// @Param id path int true "Log ID"
// @Success 200 {object} models.Log "Log details"
// @Failure 400 {object} map[string]string "Invalid log ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Log not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /logs/{id} [get]
func GetLog(c *fiber.Ctx) error {
	logID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid log ID",
		})
	}

	log, err := logService.GetLogByID(logID)
	if err != nil {
		if err.Error() == "log not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(
		models.SuccessResponse("Log retrieved successfully", log),
	)
}
