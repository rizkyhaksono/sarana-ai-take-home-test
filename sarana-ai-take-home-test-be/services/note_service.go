package services

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/constants"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/database"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/models"
)

type NoteService struct{}

func NewNoteService() *NoteService {
	return &NoteService{}
}

func (s *NoteService) CreateNote(userID uuid.UUID, title, content string) (*models.Note, error) {
	var note models.Note
	err := database.DB.QueryRow(
		"INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id, user_id, title, content, image_path, created_at, updated_at",
		userID, title, content,
	).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.ImagePath, &note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		return nil, errors.New(constants.ErrCreatingNote)
	}

	return &note, nil
}

func (s *NoteService) GetNotesByUserID(userID uuid.UUID) ([]models.Note, error) {
	rows, err := database.DB.Query(
		"SELECT id, user_id, title, content, image_path, created_at, updated_at FROM notes WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, errors.New(constants.ErrFetchingNotes)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.ImagePath, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, errors.New(constants.ErrFetchingNotes)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (s *NoteService) GetNoteByID(noteID, userID uuid.UUID) (*models.Note, error) {
	var note models.Note
	err := database.DB.QueryRow(
		"SELECT id, user_id, title, content, image_path, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2",
		noteID, userID,
	).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.ImagePath, &note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(constants.ErrNoteNotFound)
		}
		return nil, errors.New(constants.ErrFetchingNotes)
	}

	return &note, nil
}

func (s *NoteService) DeleteNote(noteID, userID uuid.UUID) error {
	result, err := database.DB.Exec(
		"DELETE FROM notes WHERE id = $1 AND user_id = $2",
		noteID, userID,
	)
	if err != nil {
		return errors.New(constants.ErrDeletingNote)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New(constants.ErrNoteNotFound)
	}

	return nil
}

func (s *NoteService) UploadImage(noteID, userID uuid.UUID, file *multipart.FileHeader) (string, error) {
	// Verify note ownership
	_, err := s.GetNoteByID(noteID, userID)
	if err != nil {
		return "", err
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !strings.Contains(constants.AllowedImageTypes, ext) {
		return "", errors.New(constants.ErrInvalidFileType)
	}

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(constants.UploadDir, 0755); err != nil {
		return "", errors.New(constants.ErrSavingFile)
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filepath := filepath.Join(constants.UploadDir, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", errors.New(constants.ErrSavingFile)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", errors.New(constants.ErrSavingFile)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return "", errors.New(constants.ErrSavingFile)
	}

	// Update note with image path
	_, err = database.DB.Exec(
		"UPDATE notes SET image_path = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		filepath, noteID,
	)
	if err != nil {
		return "", errors.New(constants.ErrUpdatingNote)
	}

	return filepath, nil
}
