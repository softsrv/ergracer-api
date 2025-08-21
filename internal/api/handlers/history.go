package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	db *sql.DB
}

func NewHistoryHandler(db *sql.DB) *HistoryHandler {
	return &HistoryHandler{db: db}
}

type RaceHistory struct {
	RaceID       int     `json:"race_id"`
	RaceUUID     string  `json:"race_uuid"`
	Distance     int     `json:"distance"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
	FinishedAt   *string `json:"finished_at"`
	UserStatus   string  `json:"user_status"`
	UserDistance int     `json:"user_distance"`
	UserPace     *string `json:"user_pace"`
	UserPosition *int    `json:"user_position"`
	Participants []RaceParticipantHistory `json:"participants"`
}

type RaceParticipantHistory struct {
	UserID   int     `json:"user_id"`
	Username string  `json:"username"`
	Status   string  `json:"status"`
	Distance int     `json:"distance"`
	Pace     *string `json:"pace"`
	Position *int    `json:"position"`
}

func (h *HistoryHandler) GetUserRaceHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	query := `
		SELECT 
			r.id, r.uuid, r.distance, r.status, r.created_at, r.finished_at,
			rp.status, rp.current_distance, rp.pace, rp.position
		FROM races r
		JOIN race_participants rp ON r.id = rp.race_id
		WHERE rp.user_id = $1
		ORDER BY r.created_at DESC`

	rows, err := h.db.Query(query, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get race history"})
		return
	}
	defer rows.Close()

	var races []RaceHistory
	for rows.Next() {
		var race RaceHistory
		err := rows.Scan(
			&race.RaceID, &race.RaceUUID, &race.Distance, &race.Status,
			&race.CreatedAt, &race.FinishedAt, &race.UserStatus,
			&race.UserDistance, &race.UserPace, &race.UserPosition,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan race history"})
			return
		}

		participants, err := h.getRaceParticipants(race.RaceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get race participants"})
			return
		}
		race.Participants = participants

		races = append(races, race)
	}

	c.JSON(http.StatusOK, gin.H{"races": races})
}

func (h *HistoryHandler) getRaceParticipants(raceID int) ([]RaceParticipantHistory, error) {
	query := `
		SELECT 
			rp.user_id, u.username, rp.status, rp.current_distance, rp.pace, rp.position
		FROM race_participants rp
		JOIN users u ON rp.user_id = u.id
		WHERE rp.race_id = $1
		ORDER BY COALESCE(rp.position, 999), rp.joined_at`

	rows, err := h.db.Query(query, raceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []RaceParticipantHistory
	for rows.Next() {
		var p RaceParticipantHistory
		err := rows.Scan(
			&p.UserID, &p.Username, &p.Status, &p.Distance, &p.Pace, &p.Position,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}

	return participants, nil
}