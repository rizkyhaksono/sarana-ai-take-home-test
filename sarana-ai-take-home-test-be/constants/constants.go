package constants

const (
	ErrInvalidRequestBody = "Invalid request body"
	ErrHashingPassword    = "Error hashing password"
	ErrCreatingUser       = "Error creating user"
	ErrEmailExists        = "Email already exists"
	ErrInvalidCredentials = "Invalid credentials"
	ErrUserNotFound       = "User not found"
	ErrGeneratingToken    = "Error generating token"
	ErrCreatingNote       = "Error creating note"
	ErrFetchingNotes      = "Error fetching notes"
	ErrNoteNotFound       = "Note not found"
	ErrUnauthorized       = "Unauthorized"
	ErrDeletingNote       = "Error deleting note"
	ErrInvalidFileType    = "Invalid file type"
	ErrSavingFile         = "Error saving file"
	ErrUpdatingNote       = "Error updating note"
	ErrInvalidNoteID      = "Invalid note ID"
)

const (
	MaxFileSize       = 5 * 1024 * 1024 // 5MB
	AllowedImageTypes = ".jpg,.jpeg,.png,.gif"
	UploadDir         = "./uploads"
)

const (
	JWTExpiration = 24 // hours
)
