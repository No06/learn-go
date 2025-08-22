package handler

import (
	"net/http"

	"hinoob.net/learn-go/internal/model"
	"hinoob.net/learn-go/internal/service"

	"github.com/gin-gonic/gin"
)

// --- Request and Response Structs ---

type CreateUserRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required,min=6"`
	FullName   string `json:"full_name"`
	Role       string `json:"role" binding:"required,oneof=student teacher"`
	TeacherIDs []uint `json:"teacher_ids"` // Only used if role is 'student'
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// --- Handlers ---

// CreateUserHandler handles the user creation request
func CreateUserHandler(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := service.CreateUser(req.Username, req.Password, req.FullName, model.Role(req.Role), req.TeacherIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": UserResponse{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Role:     string(user.Role),
		},
	})
}

// LoginHandler handles the user login request
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Role:     string(user.Role),
		},
	})
}
