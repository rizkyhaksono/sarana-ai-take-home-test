package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/constants"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/models"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/services"
)

var noteService = services.NewNoteService()

// CreateNote creates a new note for the authenticated user
func CreateNote(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("userID").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidRequestBody,
		})
	}

	var req models.CreateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidRequestBody,
		})
	}

	note, err := noteService.CreateNote(userID, req.Title, req.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(note)
}

// GetNotes retrieves all notes for the authenticated user
func GetNotes(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("userID").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.ErrInvalidRequestBody,
		})
	}

	notes, err := noteService.GetNotesByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(notes)
}

// GetNote retrieves a single note by ID
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
