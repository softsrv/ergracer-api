package services

import (
	"database/sql"
	"fmt"

	"ergracer-api/internal/models"
	"ergracer-api/internal/utils"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(email, username, password string) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken()
	if err != nil {
		return nil, err
	}

	var user models.User
	query := `
		INSERT INTO users (email, username, password_hash, email_verify_token)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, username, email_verified, created_at, updated_at`
	
	err = s.db.QueryRow(query, email, username, hashedPassword, token).Scan(
		&user.ID, &user.Email, &user.Username, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	user.EmailVerifyToken = &token
	return &user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, username, password_hash, email_verified, created_at, updated_at FROM users WHERE email = $1`
	
	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, username, email_verified, created_at, updated_at FROM users WHERE id = $1`
	
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) VerifyEmail(token string) error {
	query := `UPDATE users SET email_verified = true, email_verify_token = NULL WHERE email_verify_token = $1`
	result, err := s.db.Exec(query, token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("invalid verification token")
	}

	return nil
}

func (s *UserService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.EmailVerified {
		return nil, fmt.Errorf("email not verified")
	}

	return user, nil
}