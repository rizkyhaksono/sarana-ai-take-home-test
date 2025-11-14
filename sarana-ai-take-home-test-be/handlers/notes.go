package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/constants"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/models"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/services"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/utils"
)

var noteService = services.NewNoteService()

// CreateNote creates a new note for the authenticated user
// @Summary Create a new note
// @Description Create a new note for the authenticated user with optional image upload
// @Tags Notes
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Note title"
// @Param content formData string true "Note content"
// @Param image_path formData file false "Optional image file"
// @Success 201 {object} models.Note "Note created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notes [post]
func CreateNote(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("userID").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.ErrorResponse("INVALID_USER", "Invalid user ID", constants.ErrInvalidRequestBody),
		)
	}

	title := c.FormValue("title")
	content := c.FormValue("content")

	if title == "" || content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.ErrorResponse("INVALID_INPUT", "Title and content are required", constants.ErrInvalidRequestBody),
		)
	}

	file, err := c.FormFile("image_path")
	if err == nil && file != nil {
		note, err := noteService.CreateNoteWithImage(userID, title, content, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				models.ErrorResponse("CREATE_NOTE_ERROR", "Failed to create note", err.Error()),
			)
		}
		return c.Status(fiber.StatusCreated).JSON(
			models.SuccessResponse("Note created successfully", note),
		)
	}

	note, err := noteService.CreateNote(userID, title, content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.ErrorResponse("CREATE_NOTE_ERROR", "Failed to create note", err.Error()),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		models.SuccessResponse("Note created successfully", note),
	)
}

// GetNotes retrieves all notes for the authenticated user with search, sort, and pagination
// @Summary Get all notes
// @Description Retrieve all notes for the authenticated user with optional search, sorting, and pagination
// @Tags Notes
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search in title and content"
// @Param sort_by query string false "Sort by field (created_at, updated_at, title)" default(created_at)
// @Param order query string false "Sort order (ASC, DESC)" default(DESC)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} services.NotesResponse "Paginated list of notes"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notes [get]
func GetNotes(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("userID").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidRequestBody,
		})
	}

	params := utils.PaginationParams{
		Search: c.Query("search", ""),
		SortBy: c.Query("sort_by", "created_at"),
		Order:  c.Query("order", "DESC"),
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
	}

	result, err := noteService.GetNotesWithParams(userID, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			models.ErrorResponse("GET_NOTES_ERROR", "Failed to retrieve notes", err.Error()),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		models.SuccessResponse("Notes retrieved successfully", result),
	)
}

// GetNote retrieves a single note by ID
// @Summary Get a note by ID
// @Description Retrieve a specific note by its ID
// @Tags Notes
// @Produce json
// @Security BearerAuth
// @Param id path string true "Note ID (UUID)"
// @Success 200 {object} models.Note "Note details"
// @Failure 400 {object} map[string]string "Invalid note ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Note not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notes/{id} [get]
func GetNote(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("userID").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidRequestBody,
		})
	}

	noteID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid note ID",
		})
	}

	note, err := noteService.GetNoteByID(noteID, userID)
	if err != nil {
		if err.Error() == constants.ErrNoteNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(note)
}

// DeleteNote deletes a note by ID (with ownership check)
// @Summary Delete a note
// @Description Delete a specific note by its ID (with ownership verification)
// @Tags Notes
// @Produce json
// @Security BearerAuth
// @Param id path string true "Note ID (UUID)"
// @Success 200 {object} map[string]string "Note deleted successfully"
// @Failure 400 {object} map[string]string "Invalid note ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Note not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notes/{id} [delete]
func DeleteNote(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("userID").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidRequestBody,
		})
	}

	noteID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidNoteID,
		})
	}

	err = noteService.DeleteNote(noteID, userID)
	if err != nil {
		if err.Error() == constants.ErrNoteNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Note deleted successfully",
	})
}

// UploadNoteImage handles image upload for a specific note
// @Summary Upload an image to a note
// @Description Upload an image file to attach to a specific note
// @Tags Notes
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Note ID (UUID)"
// @Param image formData file true "Image file to upload"
// @Success 200 {object} map[string]string "Image uploaded successfully"
// @Failure 400 {object} map[string]string "Invalid request or file type"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Note not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notes/{id}/image [post]
func UploadNoteImage(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("userID").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidRequestBody,
		})
	}

	noteID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidNoteID,
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	imagePath, err := noteService.UploadImage(noteID, userID, file)
	if err != nil {
		switch err.Error() {
		case constants.ErrNoteNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		case constants.ErrInvalidFileType:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Image uploaded successfully",
		"image_path": imagePath,
	})
}

// GetNoteImage serves the image file for a specific note
// @Summary Get note image
// @Description Retrieve the image file attached to a specific note
// @Tags Notes
// @Produce image/jpeg,image/png,image/gif,image/webp
// @Security BearerAuth
// @Param id path string true "Note ID (UUID)"
// @Success 200 {file} binary "Image file"
// @Failure 400 {object} map[string]string "Invalid note ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Note or image not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notes/{id}/image [get]
func GetNoteImage(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("userID").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidRequestBody,
		})
	}

	noteID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidNoteID,
		})
	}

	// Get the note to verify ownership and get image path
	note, err := noteService.GetNoteByID(noteID, userID)
	if err != nil {
		if err.Error() == constants.ErrNoteNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if note has an image
	if note.ImagePath == nil || *note.ImagePath == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No image found for this note",
		})
	}

	// Send the file
	return c.SendFile(*note.ImagePath)
}
