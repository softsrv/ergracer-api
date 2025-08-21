package services

import (
	"database/sql"
	"fmt"
	"time"

	"ergracer-api/internal/models"

	"github.com/google/uuid"
)

type RaceService struct {
	db *sql.DB
}

func NewRaceService(db *sql.DB) *RaceService {
	return &RaceService{db: db}
}

func (s *RaceService) CreateRace(userID, distance int) (*models.Race, error) {
	raceUUID := uuid.New().String()

	var race models.Race
	query := `
		INSERT INTO races (uuid, distance, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, uuid, distance, status, created_by, created_at`
	
	err := s.db.QueryRow(query, raceUUID, distance, userID).Scan(
		&race.ID, &race.UUID, &race.Distance, &race.Status, &race.CreatedBy, &race.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(
		"INSERT INTO race_participants (race_id, user_id) VALUES ($1, $2)",
		race.ID, userID,
	)
	if err != nil {
		return nil, err
	}

	return &race, nil
}

func (s *RaceService) JoinRace(raceUUID string, userID int) error {
	var raceID int
	err := s.db.QueryRow("SELECT id FROM races WHERE uuid = $1 AND status = 'waiting'", raceUUID).Scan(&raceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("race not found or already started")
		}
		return err
	}

	_, err = s.db.Exec(
		"INSERT INTO race_participants (race_id, user_id) VALUES ($1, $2) ON CONFLICT (race_id, user_id) DO NOTHING",
		raceID, userID,
	)
	return err
}

func (s *RaceService) SetReadyStatus(raceID, userID int, ready bool) error {
	status := "not_ready"
	if ready {
		status = "ready"
	}

	_, err := s.db.Exec(
		"UPDATE race_participants SET status = $1 WHERE race_id = $2 AND user_id = $3",
		status, raceID, userID,
	)
	return err
}

func (s *RaceService) GetRaceByUUID(raceUUID string) (*models.Race, error) {
	var race models.Race
	query := `
		SELECT id, uuid, distance, status, created_by, created_at, started_at, finished_at, countdown_at
		FROM races WHERE uuid = $1`
	
	err := s.db.QueryRow(query, raceUUID).Scan(
		&race.ID, &race.UUID, &race.Distance, &race.Status, &race.CreatedBy,
		&race.CreatedAt, &race.StartedAt, &race.FinishedAt, &race.CountdownAt,
	)
	if err != nil {
		return nil, err
	}

	return &race, nil
}

func (s *RaceService) GetRaceParticipants(raceID int) ([]models.RaceParticipant, error) {
	query := `
		SELECT id, race_id, user_id, status, current_distance, finished_at, pace, position, joined_at
		FROM race_participants WHERE race_id = $1`
	
	rows, err := s.db.Query(query, raceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.RaceParticipant
	for rows.Next() {
		var p models.RaceParticipant
		err := rows.Scan(
			&p.ID, &p.RaceID, &p.UserID, &p.Status, &p.CurrentDistance,
			&p.FinishedAt, &p.Pace, &p.Position, &p.JoinedAt,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}

	return participants, nil
}

func (s *RaceService) CheckAndStartCountdown(raceID int) error {
	var totalParticipants, readyParticipants int
	
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM race_participants WHERE race_id = $1",
		raceID,
	).Scan(&totalParticipants)
	if err != nil {
		return err
	}

	err = s.db.QueryRow(
		"SELECT COUNT(*) FROM race_participants WHERE race_id = $1 AND status = 'ready'",
		raceID,
	).Scan(&readyParticipants)
	if err != nil {
		return err
	}

	if totalParticipants > 1 && totalParticipants == readyParticipants {
		countdownTime := time.Now().Add(10 * time.Second)
		_, err = s.db.Exec(
			"UPDATE races SET status = 'countdown', countdown_at = $1 WHERE id = $2",
			countdownTime, raceID,
		)
		return err
	}

	return nil
}

func (s *RaceService) StartRace(raceID int) error {
	now := time.Now()
	_, err := s.db.Exec(
		"UPDATE races SET status = 'active', started_at = $1 WHERE id = $2 AND status = 'countdown'",
		now, raceID,
	)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		"UPDATE race_participants SET status = 'racing' WHERE race_id = $1 AND status = 'ready'",
		raceID,
	)
	return err
}

func (s *RaceService) UpdateRaceProgress(raceID, userID, distance int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO race_updates (race_id, user_id, distance) VALUES ($1, $2, $3)",
		raceID, userID, distance,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE race_participants SET current_distance = $1 WHERE race_id = $2 AND user_id = $3",
		distance, raceID, userID,
	)
	if err != nil {
		return err
	}

	var raceDistance int
	err = tx.QueryRow("SELECT distance FROM races WHERE id = $1", raceID).Scan(&raceDistance)
	if err != nil {
		return err
	}

	if distance >= raceDistance {
		now := time.Now()
		_, err = tx.Exec(
			"UPDATE race_participants SET status = 'finished', finished_at = $1 WHERE race_id = $2 AND user_id = $3",
			now, raceID, userID,
		)
		if err != nil {
			return err
		}

		err = s.checkRaceCompletion(tx, raceID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *RaceService) checkRaceCompletion(tx *sql.Tx, raceID int) error {
	var totalParticipants, finishedParticipants int
	
	err := tx.QueryRow(
		"SELECT COUNT(*) FROM race_participants WHERE race_id = $1",
		raceID,
	).Scan(&totalParticipants)
	if err != nil {
		return err
	}

	err = tx.QueryRow(
		"SELECT COUNT(*) FROM race_participants WHERE race_id = $1 AND status = 'finished'",
		raceID,
	).Scan(&finishedParticipants)
	if err != nil {
		return err
	}

	if totalParticipants == finishedParticipants {
		now := time.Now()
		_, err = tx.Exec(
			"UPDATE races SET status = 'finished', finished_at = $1 WHERE id = $2",
			now, raceID,
		)
		if err != nil {
			return err
		}

		err = s.calculateRaceResults(tx, raceID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *RaceService) calculateRaceResults(tx *sql.Tx, raceID int) error {
	query := `
		WITH race_data AS (
			SELECT r.distance, r.started_at
			FROM races r
			WHERE r.id = $1
		),
		participant_times AS (
			SELECT 
				rp.id,
				rp.user_id,
				rp.finished_at,
				rd.started_at,
				rd.distance,
				EXTRACT(EPOCH FROM (rp.finished_at - rd.started_at)) as total_seconds,
				ROW_NUMBER() OVER (ORDER BY rp.finished_at) as position
			FROM race_participants rp
			CROSS JOIN race_data rd
			WHERE rp.race_id = $1 AND rp.status = 'finished'
		)
		UPDATE race_participants rp
		SET 
			pace = pt.pace_formatted,
			position = pt.position
		FROM (
			SELECT 
				pt.id,
				pt.position,
				LPAD(FLOOR(pt.total_seconds / pt.distance * 500 / 60)::text, 2, '0') || ':' || 
				LPAD(FLOOR(pt.total_seconds / pt.distance * 500 % 60)::text, 2, '0') as pace_formatted
			FROM participant_times pt
		) pt
		WHERE rp.id = pt.id`

	_, err := tx.Exec(query, raceID)
	return err
}