package models

import (
	"time"
)

type User struct {
	ID                int       `json:"id" db:"id"`
	Email             string    `json:"email" db:"email"`
	Username          string    `json:"username" db:"username"`
	PasswordHash      string    `json:"-" db:"password_hash"`
	EmailVerified     bool      `json:"email_verified" db:"email_verified"`
	EmailVerifyToken  *string   `json:"-" db:"email_verify_token"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type Session struct {
	ID                 int       `json:"id" db:"id"`
	UserID             int       `json:"user_id" db:"user_id"`
	RefreshTokenHash   string    `json:"-" db:"refresh_token_hash"`
	DeviceType         string    `json:"device_type" db:"device_type"`
	UserAgent          string    `json:"user_agent" db:"user_agent"`
	IPAddress          string    `json:"ip_address" db:"ip_address"`
	ExpiresAt          time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type Friendship struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	FriendID   int       `json:"friend_id" db:"friend_id"`
	Status     string    `json:"status" db:"status"` // pending, accepted
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	AcceptedAt *time.Time `json:"accepted_at" db:"accepted_at"`
}