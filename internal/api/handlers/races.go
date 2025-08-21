package handlers

import (
	"net/http"
	"strconv"

	"ergracer-api/internal/services"

	"github.com/gin-gonic/gin"
)

type RacesHandler struct {
	raceService *services.RaceService
}

func NewRacesHandler(raceService *services.RaceService) *RacesHandler {
	return &RacesHandler{
		raceService: raceService,
	}
}

type CreateRaceRequest struct {
	Distance int `json:"distance" binding:"required,min=100"`
}

type JoinRaceRequest struct {
	RaceUUID string `json:"race_uuid" binding:"required"`
}

type SetReadyRequest struct {
	Ready bool `json:"ready"`
}

type UpdateProgressRequest struct {
	Distance int `json:"distance" binding:"required,min=0"`
}

func (h *RacesHandler) CreateRace(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateRaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	race, err := h.raceService.CreateRace(userID.(int), req.Distance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create race"})
		return
	}

	c.JSON(http.StatusCreated, race)
}

func (h *RacesHandler) JoinRace(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req JoinRaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.raceService.JoinRace(req.RaceUUID, userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined race successfully"})
}

func (h *RacesHandler) GetRace(c *gin.Context) {
	raceUUID := c.Param("uuid")
	
	race, err := h.raceService.GetRaceByUUID(raceUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Race not found"})
		return
	}

	participants, err := h.raceService.GetRaceParticipants(race.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get participants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"race":         race,
		"participants": participants,
	})
}

func (h *RacesHandler) SetReady(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	raceIDStr := c.Param("raceId")
	raceID, err := strconv.Atoi(raceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid race ID"})
		return
	}

	var req SetReadyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.raceService.SetReadyStatus(raceID, userID.(int), req.Ready)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to set ready status"})
		return
	}

	if req.Ready {
		err = h.raceService.CheckAndStartCountdown(raceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check countdown"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ready status updated"})
}

func (h *RacesHandler) UpdateProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	raceIDStr := c.Param("raceId")
	raceID, err := strconv.Atoi(raceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid race ID"})
		return
	}

	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.raceService.UpdateRaceProgress(raceID, userID.(int), req.Distance)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Progress updated"})
}

func (h *RacesHandler) StartRace(c *gin.Context) {
	raceIDStr := c.Param("raceId")
	raceID, err := strconv.Atoi(raceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid race ID"})
		return
	}

	err = h.raceService.StartRace(raceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to start race"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Race started"})
}