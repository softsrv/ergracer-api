package handlers

import (
	"net/http"
	"strconv"

	"ergracer-api/internal/services"

	"github.com/gin-gonic/gin"
)

type FriendsHandler struct {
	friendshipService *services.FriendshipService
	userService       *services.UserService
}

func NewFriendsHandler(friendshipService *services.FriendshipService, userService *services.UserService) *FriendsHandler {
	return &FriendsHandler{
		friendshipService: friendshipService,
		userService:       userService,
	}
}

type InviteFriendRequest struct {
	FriendID int `json:"friend_id" binding:"required"`
}

func (h *FriendsHandler) InviteFriend(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req InviteFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userID.(int) == req.FriendID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot invite yourself"})
		return
	}

	_, err := h.userService.GetUserByID(req.FriendID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = h.friendshipService.InviteFriend(userID.(int), req.FriendID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend invitation sent"})
}

func (h *FriendsHandler) AcceptFriendship(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	friendIDStr := c.Param("friendId")
	friendID, err := strconv.Atoi(friendIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	err = h.friendshipService.AcceptFriendship(userID.(int), friendID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friendship accepted"})
}

func (h *FriendsHandler) GetFriends(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	friends, err := h.friendshipService.GetFriends(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get friends"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"friends": friends})
}

func (h *FriendsHandler) GetPendingInvitations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	invitations, err := h.friendshipService.GetPendingInvitations(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending invitations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invitations": invitations})
}