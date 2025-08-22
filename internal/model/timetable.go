package model

import (
	"gorm.io/gorm"
)

// Class represents a group of students.
type Class struct {
	gorm.Model
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Students  []*User   `gorm:"many2many:class_students;"`
	Courses   []Course  // A class can have many courses scheduled
}

// TimeSlot defines a fixed period in the school day, e.g., "1st Period".
// These are managed by an admin in the backend.
type TimeSlot struct {
	gorm.Model
	Name      string `gorm:"type:varchar(50);not null"` // e.g., "Period 1", "午休"
	StartTime string `gorm:"type:varchar(5);not null"`  // Format: HH:MM
	EndTime   string `gorm:"type:varchar(5);not null"`  // Format: HH:MM
}

// Course represents a specific subject taught to a specific class at a specific time.
type Course struct {
	gorm.Model
	SubjectName string    `gorm:"type:varchar(100);not null"` // e.g., "Mathematics", "History"
	ClassID     uint      `gorm:"not null"`
	TeacherID   uint      `gorm:"not null"`
	TimeSlotID  uint      `gorm:"not null"`
	DayOfWeek   int       `gorm:"not null"` // 1 for Monday, 2 for Tuesday, etc.
	
	Class       Class     `gorm:"foreignKey:ClassID"`
	Teacher     User      `gorm:"foreignKey:TeacherID"`
	TimeSlot    TimeSlot  `gorm:"foreignKey:TimeSlotID"`
}