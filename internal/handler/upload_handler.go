package handler

import (
	"net/http"

	"hinoob.net/learn-go/pkg/oss"

	"github.com/gin-gonic/gin"
)

// UploadFileHandler handles file uploads from a multipart form.
func UploadFileHandler(c *gin.Context) {
	// "file" is the name of the form field
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed: " + err.Error()})
		return
	}

	// Upload the file to OSS
	fileURL, err := oss.UploadFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to cloud storage: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"file_url": fileURL,
	})
}
