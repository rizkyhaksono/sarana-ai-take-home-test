package services

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/rizkyhaksono/sarana-ai-take-home-test/constants"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/database"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/models"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(email, password string) (*models.User, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", errors.New(constants.ErrHashingPassword)
	}

	var user models.User
	err = database.DB.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, created_at",
		email, string(hashedPassword),
	).Scan(&user.ID, &user.Email, &user.CreatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, "", errors.New(constants.ErrEmailExists)
		}
		return nil, "", errors.New(constants.ErrCreatingUser)
	}

	token, err := utils.GenerateJWT(user.ID.String(), user.Email)
	if err != nil {
		return nil, "", errors.New(constants.ErrGeneratingToken)
	}

	return &user, token, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	var user models.User
	var hashedPassword string

	err := database.DB.QueryRow(
		"SELECT id, email, password, created_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &hashedPassword, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.New(constants.ErrInvalidCredentials)
		}
		return nil, "", errors.New(constants.ErrUserNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil, "", errors.New(constants.ErrInvalidCredentials)
	}

	token, err := utils.GenerateJWT(user.ID.String(), user.Email)
	if err != nil {
		return nil, "", errors.New(constants.ErrGeneratingToken)
	}

	return &user, token, nil
}
