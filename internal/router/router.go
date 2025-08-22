package router

import (
	"hinoob.net/learn-go/internal/handler"
	"hinoob.net/learn-go/internal/middleware"
	"hinoob.net/learn-go/internal/model"
	"hinoob.net/learn-go/internal/pkg/websocket"
	"hinoob.net/learn-go/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the routes for the application
func SetupRouter(hub *websocket.Hub, liveService *service.LiveService) *gin.Engine {
	router := gin.Default()

	// Health check and WebSocket endpoints
	router.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
	router.GET("/ws", func(c *gin.Context) { handler.ServeWs(hub, c) })

	// API v1 group
	apiV1 := router.Group("/api/v1")
	{
		// --- Public Routes ---
		userRoutes := apiV1.Group("/users")
		{
			userRoutes.POST("/register", handler.CreateUserHandler)
			userRoutes.POST("/login", handler.LoginHandler)
		}

		// --- Authenticated Routes ---
		authGroup := apiV1.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			// Initialize handlers that depend on services
			liveHandler := handler.NewLiveHandler(liveService)

			// Upload route
			authGroup.POST("/upload", handler.UploadFileHandler)

			// Message history route
			authGroup.GET("/messages/history/:userId", handler.GetMessageHistoryHandler)

			// Assignment routes
			assignmentRoutes := authGroup.Group("/assignments")
			{
				assignmentRoutes.POST("", middleware.RoleAuthMiddleware(string(model.TeacherRole)), handler.CreateAssignmentHandler)
				assignmentRoutes.GET("", handler.GetAssignmentsHandler)
				assignmentSpecificRoutes := assignmentRoutes.Group("/:assignmentId")
				{
					assignmentSpecificRoutes.POST("/submit", middleware.RoleAuthMiddleware(string(model.StudentRole)), handler.SubmitAssignmentHandler)
					assignmentSpecificRoutes.GET("/submissions", middleware.RoleAuthMiddleware(string(model.TeacherRole)), handler.GetSubmissionsForAssignmentHandler)
				}
			}

			// Timetable and Course routes
			timetableRoutes := authGroup.Group("/timetable")
			{
				timetableRoutes.POST("/slots", middleware.RoleAuthMiddleware(string(model.TeacherRole)), handler.CreateTimeSlotHandler)
				timetableRoutes.POST("/courses", middleware.RoleAuthMiddleware(string(model.TeacherRole)), handler.CreateCourseHandler)
				timetableRoutes.GET("", handler.GetTimetableHandler)
			}

			// Live Stream routes
			liveRoutes := authGroup.Group("/live")
			{
				liveRoutes.POST("/:courseId/start", middleware.RoleAuthMiddleware(string(model.TeacherRole)), liveHandler.StartStreamHandler)
				liveRoutes.POST("/:courseId/end", middleware.RoleAuthMiddleware(string(model.TeacherRole)), liveHandler.EndStreamHandler)
			}

			// Submission routes (for grading)
			submissionRoutes := authGroup.Group("/submissions")
			{
				submissionRoutes.POST("/:submissionId/grade", middleware.RoleAuthMiddleware(string(model.TeacherRole)), handler.GradeSubmissionHandler)
			}
		}
	}

	return router
}
