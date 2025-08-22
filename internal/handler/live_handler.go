package handler

import (
	"net/http"
	"strconv"

	"hinoob.net/learn-go/internal/middleware"
	"hinoob.net/learn-go/internal/service"

	"github.com/gin-gonic/gin"
)

type LiveHandler struct {
	liveService *service.LiveService
}

func NewLiveHandler(liveService *service.LiveService) *LiveHandler {
	return &LiveHandler{liveService: liveService}
}

// StartStreamHandler handles a teacher's request to start a live stream.
func (h *LiveHandler) StartStreamHandler(c *gin.Context) {
	courseID, _ := strconv.ParseUint(c.Param("courseId"), 10, 32)
	teacherID := c.GetUint(middleware.ContextUserIDKey)

	err := h.liveService.StartStreamForCourse(uint(courseID), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Live stream started successfully"})
}

// EndStreamHandler handles a teacher's request to end a live stream.
func (h *LiveHandler) EndStreamHandler(c *gin.Context) {
	courseID, _ := strconv.ParseUint(c.Param("courseId"), 10, 32)
	teacherID := c.GetUint(middleware.ContextUserIDKey)

	err := h.liveService.EndStreamForCourse(uint(courseID), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Live stream ended successfully"})
}
