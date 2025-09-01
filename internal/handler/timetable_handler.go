package handler

import (
	"net/http"

	"hinoob.net/learn-go/internal/database"
	"hinoob.net/learn-go/internal/middleware"
	"hinoob.net/learn-go/internal/model"
	"hinoob.net/learn-go/internal/service"

	"github.com/gin-gonic/gin"
)

// --- Request Structs ---

type CreateTimeSlotRequest struct {
	Name      string `json:"name" binding:"required"`
	StartTime string `json:"start_time" binding:"required,len=5"` // HH:MM
	EndTime   string `json:"end_time" binding:"required,len=5"`   // HH:MM
}

type CreateCourseRequest struct {
	SubjectName string `json:"subject_name" binding:"required"`
	ClassID     uint   `json:"class_id" binding:"required"`
	TimeSlotID  uint   `json:"time_slot_id" binding:"required"`
	DayOfWeek   int    `json:"day_of_week" binding:"required,min=1,max=7"`
}

// --- Handlers ---

// CreateTimeSlotHandler handles creating a new time slot (admin).
func CreateTimeSlotHandler(c *gin.Context) {
	var req CreateTimeSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeSlot, err := service.CreateTimeSlot(req.Name, req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create time slot"})
		return
	}

	c.JSON(http.StatusCreated, timeSlot)
}

// CreateCourseHandler handles a teacher creating a course.
func CreateCourseHandler(c *gin.Context) {
	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherID := c.GetUint(middleware.ContextUserIDKey)
	course, err := service.CreateCourseForTeacher(req.SubjectName, req.ClassID, teacherID, req.TimeSlotID, req.DayOfWeek)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// GetTimetableHandler handles fetching the timetable.
// For teachers, it returns their schedule. For students, it returns their class's schedule.
func GetTimetableHandler(c *gin.Context) {
	role := c.GetString(middleware.ContextUserRoleKey)
	userID := c.GetUint(middleware.ContextUserIDKey)

	var courses []model.Course
	var err error

	if role == string(model.TeacherRole) {
		courses, err = database.GetCoursesByTeacher(userID)
	} else {
		// This requires getting the student's class. This is a missing piece.
		// We need to be able to find which class a student belongs to.
		// For now, we will leave this part unimplemented.
		// student, err := database.GetUserWithClass(userID)
		// if err == nil {
		// 	courses, err = database.GetCoursesByClass(student.ClassID)
		// }
		courses = make([]model.Course, 0)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve timetable"})
		return
	}

	c.JSON(http.StatusOK, courses)
}
