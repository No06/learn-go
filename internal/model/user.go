package model

import (
	"gorm.io/gorm"
	"time"
)

// Role defines the user role enum
type Role string

const (
	StudentRole Role = "student"
	TeacherRole Role = "teacher"
)

// User represents a user account in the system
type User struct {
	gorm.Model
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	Role         Role      `gorm:"type:varchar(10);not null"`
	FullName     string    `gorm:"type:varchar(100)"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	// Relationships
	// For Students: The teachers they are assigned to.
	Teachers []*User `gorm:"many2many:student_teachers;"`
}

// StudentTeacher is the join table for the many-to-many relationship
// between students and teachers. GORM will create this automatically.
// We can define it explicitly if we need to add more columns to the join table.
// type StudentTeacher struct {
//    StudentID uint `gorm:"primaryKey"`
//    TeacherID uint `gorm:"primaryKey"`
// }

func (u *User) IsStudent() bool {
	return u.Role == StudentRole
}

func (u *User) IsTeacher() bool {
	return u.Role == TeacherRole
}
