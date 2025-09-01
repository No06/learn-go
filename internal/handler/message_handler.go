package handler

import (
	"net/http"
	"strconv"

	"hinoob.net/learn-go/internal/database"
	"hinoob.net/learn-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// GetMessageHistoryHandler retrieves the chat history between the logged-in user and another user.
func GetMessageHistoryHandler(c *gin.Context) {
	// 1. Get the ID of the other user from the URL parameter
	otherUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 2. Get the logged-in user's ID from the context
	currentUserID := c.GetUint(middleware.ContextUserIDKey)

	// 3. Fetch the messages from the database
	messages, err := database.GetMessagesBetweenUsers(currentUserID, uint(otherUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve message history"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
