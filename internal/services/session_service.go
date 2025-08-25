package services

import (
	"database/sql"
	"time"

	"ergracer-api/internal/models"
	"ergracer-api/internal/utils"
)

type SessionService struct {
	db *sql.DB
}

func NewSessionService(db *sql.DB) *SessionService {
	return &SessionService{db: db}
}

func (s *SessionService) CreateSession(userID int, deviceType, userAgent, ipAddress string) (string, error) {
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	refreshTokenHash, err := utils.HashRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	query := `
		INSERT INTO sessions (user_id, refresh_token_hash, device_type, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = s.db.Exec(query, userID, refreshTokenHash, deviceType, userAgent, ipAddress, expiresAt)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (s *SessionService) ValidateRefreshToken(refreshToken string) (*models.Session, error) {
	refreshTokenHash, err := utils.HashRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, user_id, refresh_token_hash, device_type, user_agent, ip_address, expires_at, created_at, updated_at
		FROM sessions
		WHERE refresh_token_hash = $1 AND expires_at > NOW()
	`

	var session models.Session
	err = s.db.QueryRow(query, refreshTokenHash).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshTokenHash,
		&session.DeviceType,
		&session.UserAgent,
		&session.IPAddress,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionService) UpdateSession(sessionID int, newRefreshToken string) error {
	refreshTokenHash, err := utils.HashRefreshToken(newRefreshToken)
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	query := `
		UPDATE sessions
		SET refresh_token_hash = $1, expires_at = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err = s.db.Exec(query, refreshTokenHash, expiresAt, sessionID)
	return err
}

func (s *SessionService) DeleteExpiredSessions() error {
	query := `DELETE FROM sessions WHERE expires_at <= NOW()`
	_, err := s.db.Exec(query)
	return err
}

func (s *SessionService) DeleteSession(refreshToken string) error {
	refreshTokenHash, err := utils.HashRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	query := `DELETE FROM sessions WHERE refresh_token_hash = $1`
	_, err = s.db.Exec(query, refreshTokenHash)
	return err
}

func (s *SessionService) DeleteUserSessions(userID int) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := s.db.Exec(query, userID)
	return err
}
