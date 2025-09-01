package service

import (
	"hinoob.net/learn-go/internal/database"
	"hinoob.net/learn-go/internal/model"
)

// --- TimeSlot Services ---

// CreateTimeSlot handles the logic for creating a new time slot.
// This would typically be an admin-only operation.
func CreateTimeSlot(name, startTime, endTime string) (*model.TimeSlot, error) {
	timeSlot := &model.TimeSlot{
		Name:      name,
		StartTime: startTime,
		EndTime:   endTime,
	}
	err := database.CreateTimeSlot(timeSlot)
	return timeSlot, err
}

// --- Course Services ---

// CreateCourseForTeacher handles the logic for a teacher creating a new course in the timetable.
func CreateCourseForTeacher(subject string, classID, teacherID, timeSlotID uint, dayOfWeek int) (*model.Course, error) {
	course := &model.Course{
		SubjectName: subject,
		ClassID:     classID,
		TeacherID:   teacherID,
		TimeSlotID:  timeSlotID,
		DayOfWeek:   dayOfWeek,
	}
	err := database.CreateCourse(course)
	return course, err
}
