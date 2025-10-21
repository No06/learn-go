package response

import "github.com/gin-gonic/gin"

// Success writes a JSON success response with payload.
func Success(ctx *gin.Context, status int, payload interface{}) {
	ctx.JSON(status, gin.H{"success": true, "data": payload})
}

// Error writes an error response with message and optional details.
func Error(ctx *gin.Context, status int, message string, details interface{}) {
	ctx.JSON(status, gin.H{
		"success": false,
		"error": gin.H{
			"message": message,
			"details": details,
		},
	})
}
