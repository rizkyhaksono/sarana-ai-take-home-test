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
	"github.com/rizkyhaksono/sarana-ai-take-home-test/utils"
)

type NoteService struct{}

func NewNoteService() *NoteService {
	return &NoteService{}
}

func (s *NoteService) CreateNote(userID uuid.UUID, title, content string) (*models.Note, error) {
	var note models.Note
	var imagePath sql.NullString
	err := database.DB.QueryRow(
		"INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id, user_id, title, content, image_path, created_at, updated_at",
		userID, title, content,
	).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &imagePath, &note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("%s: %v", constants.ErrCreatingNote, err)
	}

	if imagePath.Valid {
		note.ImagePath = &imagePath.String
	}

	return &note, nil
}

func (s *NoteService) CreateNoteWithImage(userID uuid.UUID, title, content string, file *multipart.FileHeader) (*models.Note, error) {
	note, err := s.CreateNote(userID, title, content)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !strings.Contains(constants.AllowedImageTypes, ext) {
		_ = s.DeleteNote(note.ID, userID)
		return nil, errors.New(constants.ErrInvalidFileType)
	}

	if err := os.MkdirAll(constants.UploadDir, 0755); err != nil {
		_ = s.DeleteNote(note.ID, userID)
		return nil, errors.New(constants.ErrSavingFile)
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(constants.UploadDir, filename)

	src, err := file.Open()
	if err != nil {
		_ = s.DeleteNote(note.ID, userID)
		return nil, errors.New(constants.ErrSavingFile)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		_ = s.DeleteNote(note.ID, userID)
		return nil, errors.New(constants.ErrSavingFile)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		_ = s.DeleteNote(note.ID, userID)
		return nil, errors.New(constants.ErrSavingFile)
	}

	_, err = database.DB.Exec(
		"UPDATE notes SET image_path = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		filePath, note.ID,
	)
	if err != nil {
		_ = os.Remove(filePath)
		_ = s.DeleteNote(note.ID, userID)
		return nil, errors.New(constants.ErrUpdatingNote)
	}

	note.ImagePath = &filePath

	return note, nil
}

func (s *NoteService) UpdateNote(noteID, userID uuid.UUID, title, content string) (*models.Note, error) {
	// Verify ownership first
	existingNote, err := s.GetNoteByID(noteID, userID)
	if err != nil {
		return nil, err
	}

	var imagePath sql.NullString
	err = database.DB.QueryRow(
		"UPDATE notes SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 AND user_id = $4 RETURNING id, user_id, title, content, image_path, created_at, updated_at",
		title, content, noteID, userID,
	).Scan(&existingNote.ID, &existingNote.UserID, &existingNote.Title, &existingNote.Content, &imagePath, &existingNote.CreatedAt, &existingNote.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(constants.ErrNoteNotFound)
		}
		return nil, fmt.Errorf("%s: %v", constants.ErrUpdatingNote, err)
	}

	if imagePath.Valid {
		existingNote.ImagePath = &imagePath.String
	}

	return existingNote, nil
}

func (s *NoteService) UpdateNoteWithImage(noteID, userID uuid.UUID, title, content string, file *multipart.FileHeader) (*models.Note, error) {
	// Verify ownership first
	existingNote, err := s.GetNoteByID(noteID, userID)
	if err != nil {
		return nil, err
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !strings.Contains(constants.AllowedImageTypes, ext) {
		return nil, errors.New(constants.ErrInvalidFileType)
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(constants.UploadDir, 0755); err != nil {
		return nil, errors.New(constants.ErrSavingFile)
	}

	// Generate new filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(constants.UploadDir, filename)

	// Save the file
	src, err := file.Open()
	if err != nil {
		return nil, errors.New(constants.ErrSavingFile)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, errors.New(constants.ErrSavingFile)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, errors.New(constants.ErrSavingFile)
	}

	// Update database with new image path
	var imagePath sql.NullString
	err = database.DB.QueryRow(
		"UPDATE notes SET title = $1, content = $2, image_path = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4 AND user_id = $5 RETURNING id, user_id, title, content, image_path, created_at, updated_at",
		title, content, filePath, noteID, userID,
	).Scan(&existingNote.ID, &existingNote.UserID, &existingNote.Title, &existingNote.Content, &imagePath, &existingNote.CreatedAt, &existingNote.UpdatedAt)

	if err != nil {
		// Clean up uploaded file if database update fails
		_ = os.Remove(filePath)
		if err == sql.ErrNoRows {
			return nil, errors.New(constants.ErrNoteNotFound)
		}
		return nil, fmt.Errorf("%s: %v", constants.ErrUpdatingNote, err)
	}

	// Remove old image if it exists
	if existingNote.ImagePath != nil && *existingNote.ImagePath != "" && *existingNote.ImagePath != filePath {
		_ = os.Remove(*existingNote.ImagePath)
	}

	if imagePath.Valid {
		existingNote.ImagePath = &imagePath.String
	}

	return existingNote, nil
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
		var imagePath sql.NullString
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &imagePath, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, errors.New(constants.ErrFetchingNotes)
		}
		if imagePath.Valid {
			note.ImagePath = &imagePath.String
		}
		notes = append(notes, note)
	}

	return notes, nil
}

type NotesResponse struct {
	Notes                    []models.Note `json:"notes"`
	utils.PaginationResponse `json:",inline"`
}

func (s *NoteService) GetNotesWithParams(userID uuid.UUID, params utils.PaginationParams) (*NotesResponse, error) {
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"title":      true,
	}
	utils.ValidatePaginationParams(&params, validSortFields, "created_at")

	baseQuery := "SELECT id, user_id, title, content, image_path, created_at, updated_at FROM notes"
	countQuery := "SELECT COUNT(*) FROM notes"
	whereCondition := "user_id = $1"
	searchFields := []string{"title", "content"}
	baseArgs := []interface{}{userID}

	query, countQueryFinal, args, err := utils.BuildPaginatedQuery(
		baseQuery,
		countQuery,
		whereCondition,
		searchFields,
		params,
		baseArgs,
	)
	if err != nil {
		return nil, errors.New(constants.ErrFetchingNotes)
	}

	countArgs := args[:len(args)-2]
	total, err := utils.GetTotalCount(database.DB, countQueryFinal, countArgs)
	if err != nil {
		return nil, errors.New(constants.ErrFetchingNotes)
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, errors.New(constants.ErrFetchingNotes)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		var imagePath sql.NullString
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &imagePath, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, errors.New(constants.ErrFetchingNotes)
		}
		if imagePath.Valid {
			note.ImagePath = &imagePath.String
		}
		notes = append(notes, note)
	}

	paginationMeta := utils.CalculatePaginationMetadata(total, params.Page, params.Limit)

	return &NotesResponse{
		Notes:              notes,
		PaginationResponse: paginationMeta,
	}, nil
}

func (s *NoteService) GetNoteByID(noteID, userID uuid.UUID) (*models.Note, error) {
	var note models.Note
	var imagePath sql.NullString
	err := database.DB.QueryRow(
		"SELECT id, user_id, title, content, image_path, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2",
		noteID, userID,
	).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &imagePath, &note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(constants.ErrNoteNotFound)
		}
		return nil, errors.New(constants.ErrFetchingNotes)
	}

	if imagePath.Valid {
		note.ImagePath = &imagePath.String
	}

	return &note, nil
}

func (s *NoteService) DeleteNote(noteID, userID uuid.UUID) error {
	note, err := s.GetNoteByID(noteID, userID)
	if err != nil {
		return err
	}

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

	if note.ImagePath != nil && *note.ImagePath != "" {
		_ = os.Remove(*note.ImagePath) // Ignore error if file doesn't exist
	}

	return nil
}

func (s *NoteService) UploadImage(noteID, userID uuid.UUID, file *multipart.FileHeader) (string, error) {
	note, err := s.GetNoteByID(noteID, userID)
	if err != nil {
		return "", err
	}

	if note.ImagePath != nil && *note.ImagePath != "" {
		_ = os.Remove(*note.ImagePath) // Ignore error if file doesn't exist
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !strings.Contains(constants.AllowedImageTypes, ext) {
		return "", errors.New(constants.ErrInvalidFileType)
	}

	if err := os.MkdirAll(constants.UploadDir, 0755); err != nil {
		return "", errors.New(constants.ErrSavingFile)
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filepath := filepath.Join(constants.UploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return "", errors.New(constants.ErrSavingFile)
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return "", errors.New(constants.ErrSavingFile)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", errors.New(constants.ErrSavingFile)
	}

	_, err = database.DB.Exec(
		"UPDATE notes SET image_path = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		filepath, noteID,
	)
	if err != nil {
		return "", errors.New(constants.ErrUpdatingNote)
	}

	return filepath, nil
}
