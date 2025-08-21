package models

import (
	"time"
)

type Race struct {
	ID            int       `json:"id" db:"id"`
	UUID          string    `json:"uuid" db:"uuid"`
	Distance      int       `json:"distance" db:"distance"` // meters
	Status        string    `json:"status" db:"status"`     // waiting, ready, countdown, active, finished
	CreatedBy     int       `json:"created_by" db:"created_by"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	StartedAt     *time.Time `json:"started_at" db:"started_at"`
	FinishedAt    *time.Time `json:"finished_at" db:"finished_at"`
	CountdownAt   *time.Time `json:"countdown_at" db:"countdown_at"`
}

type RaceParticipant struct {
	ID             int       `json:"id" db:"id"`
	RaceID         int       `json:"race_id" db:"race_id"`
	UserID         int       `json:"user_id" db:"user_id"`
	Status         string    `json:"status" db:"status"`         // not_ready, ready, racing, finished
	CurrentDistance int      `json:"current_distance" db:"current_distance"` // meters
	FinishedAt     *time.Time `json:"finished_at" db:"finished_at"`
	Pace           *string   `json:"pace" db:"pace"`             // mm:ss per 500m, calculated when finished
	Position       *int      `json:"position" db:"position"`     // 1st, 2nd, 3rd, etc.
	JoinedAt       time.Time `json:"joined_at" db:"joined_at"`
}

type RaceUpdate struct {
	ID        int       `json:"id" db:"id"`
	RaceID    int       `json:"race_id" db:"race_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Distance  int       `json:"distance" db:"distance"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}