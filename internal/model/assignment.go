package model

import (
	"gorm.io/gorm"
	"time"
)

// Assignment represents a homework assignment created by a teacher.
type Assignment struct {
	gorm.Model
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	DueDate     time.Time
	CreatedByID uint      `gorm:"not null"` // Teacher's User ID
	CreatedBy   User      `gorm:"foreignKey:CreatedByID"`
}

// Submission represents a student's submission for an assignment.
type Submission struct {
	gorm.Model
	AssignmentID uint   `gorm:"not null"`
	StudentID    uint   `gorm:"not null"`
	Content      string `gorm:"type:text"`
	Grade        string `gorm:"type:varchar(50)"` // e.g., "A+", "85/100", "Excellent"
	
	Assignment   Assignment `gorm:"foreignKey:AssignmentID"`
	Student      User       `gorm:"foreignKey:StudentID"`
	Comments     []Comment  `gorm:"foreignKey:SubmissionID"`
	Files        []SubmissionFile `gorm:"foreignKey:SubmissionID"`
}

// Comment is a message from a teacher on a submission.
type Comment struct {
	gorm.Model
	SubmissionID uint   `gorm:"not null"`
	UserID       uint   `gorm:"not null"` // Teacher's ID
	Content      string `gorm:"type:text"`
	IsImage      bool   `gorm:"default:false"` // True if Content is a URL to an image
	
	User         User   `gorm:"foreignKey:UserID"`
}

// SubmissionFile holds URLs for files uploaded by a student for a submission.
type SubmissionFile struct {
	gorm.Model
	SubmissionID uint   `gorm:"not null"`
	FileURL      string `gorm:"type:varchar(255);not null"`
}
