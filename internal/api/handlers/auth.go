package handlers

import (
	"database/sql"
	"net/http"

	"ergracer-api/internal/config"
	"ergracer-api/internal/services"
	"ergracer-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService    *services.UserService
	sessionService *services.SessionService
	config         *config.Config
}

func NewAuthHandler(userService *services.UserService, sessionService *services.SessionService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		sessionService: sessionService,
		config:         config,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(req.Email, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists or invalid data"})
		return
	}

	if user.EmailVerifyToken != nil {
		err = utils.SendVerificationEmail(
			user.Email,
			*user.EmailVerifyToken,
			h.config.AppURL(),
			h.config.MailgunDomain(),
			h.config.MailgunAPIKey(),
			h.config.MailgunFromEmail(),
			h.config.MailgunFromName(),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully. Please check your email for verification.",
		"user":    user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := utils.GenerateJWT(user.ID, h.config.JWTSecret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	deviceType := utils.DetectDeviceType(c.GetHeader("User-Agent"))
	refreshToken, err := h.sessionService.CreateSession(
		user.ID,
		deviceType,
		c.GetHeader("User-Agent"),
		c.ClientIP(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
		return
	}

	err := h.userService.VerifyEmail(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.sessionService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate refresh token"})
		}
		return
	}

	newAccessToken, err := utils.GenerateJWT(session.UserID, h.config.JWTSecret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	err = h.sessionService.UpdateSession(session.ID, newRefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}