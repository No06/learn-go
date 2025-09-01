package handler

import (
	"net/http"
	"strconv"
	"time"

	"hinoob.net/learn-go/internal/database"
	"hinoob.net/learn-go/internal/middleware"
	"hinoob.net/learn-go/internal/model"
	"hinoob.net/learn-go/internal/service"

	"github.com/gin-gonic/gin"
)

// --- Request Structs ---

type CreateAssignmentRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date" binding:"required"` // e.g., "2024-12-31T23:59:59Z"
}

type SubmitAssignmentRequest struct {
	Content  string   `json:"content"`
	FileURLs []string `json:"file_urls"`
}

type GradeSubmissionRequest struct {
	Grade       string `json:"grade" binding:"required"`
	CommentText string `json:"comment_text"`
	IsImage     bool   `json:"is_image"`
}

// --- Handlers ---

// CreateAssignmentHandler handles requests from teachers to create assignments.
func CreateAssignmentHandler(c *gin.Context) {
	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date format. Use RFC3339 (e.g., 2024-12-31T23:59:59Z)"})
		return
	}

	teacherID := c.GetUint(middleware.ContextUserIDKey)
	assignment, err := service.CreateAssignmentForTeacher(req.Title, req.Description, dueDate, teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assignment"})
		return
	}

	c.JSON(http.StatusCreated, assignment)
}

// GetAssignmentsHandler handles fetching assignments.
// For teachers, it returns assignments they created.
// For students, it should return assignments for their classes (TODO).
func GetAssignmentsHandler(c *gin.Context) {
	role := c.GetString(middleware.ContextUserRoleKey)
	userID := c.GetUint(middleware.ContextUserIDKey)

	var assignments []model.Assignment
	var err error

	if role == string(model.TeacherRole) {
		assignments, err = database.GetAssignmentsByTeacherID(userID)
	} else {
		// TODO: Implement logic for students to get assignments
		// This would involve checking their class/teacher associations.
		// For now, return an empty list for students.
		assignments = make([]model.Assignment, 0)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assignments"})
		return
	}

	c.JSON(http.StatusOK, assignments)
}

// SubmitAssignmentHandler handles a student submitting their work.
func SubmitAssignmentHandler(c *gin.Context) {
	var req SubmitAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignmentID, _ := strconv.ParseUint(c.Param("assignmentId"), 10, 32)
	studentID := c.GetUint(middleware.ContextUserIDKey)

	submission, err := service.CreateOrUpdateSubmission(uint(assignmentID), studentID, req.Content, req.FileURLs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create submission"})
		return
	}

	c.JSON(http.StatusCreated, submission)
}

// GetSubmissionsForAssignmentHandler handles a teacher fetching all submissions for an assignment.
func GetSubmissionsForAssignmentHandler(c *gin.Context) {
	assignmentID, _ := strconv.ParseUint(c.Param("assignmentId"), 10, 32)

	// Optional: Add a check to ensure the requesting teacher owns the assignment

	submissions, err := database.GetSubmissionsByAssignmentID(uint(assignmentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submissions"})
		return
	}

	c.JSON(http.StatusOK, submissions)
}

// GradeSubmissionHandler handles a teacher grading a single submission.
func GradeSubmissionHandler(c *gin.Context) {
	var req GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	submissionID, _ := strconv.ParseUint(c.Param("submissionId"), 10, 32)
	teacherID := c.GetUint(middleware.ContextUserIDKey)

	submission, err := service.GradeAndCommentOnSubmission(uint(submissionID), teacherID, req.Grade, req.CommentText, req.IsImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to grade submission"})
		return
	}

	c.JSON(http.StatusOK, submission)
}
