package repository

import "hinoob.net/learn-go/internal/model"

// --- Class Repository ---

func CreateClass(class *model.Class) error {
	return DB.Create(class).Error
}

func GetClassByID(id uint) (*model.Class, error) {
	var class model.Class
	err := DB.Preload("Students").First(&class, id).Error
	return &class, err
}

// --- TimeSlot Repository ---

func CreateTimeSlot(timeSlot *model.TimeSlot) error {
	return DB.Create(timeSlot).Error
}

// --- Course Repository ---

func CreateCourse(course *model.Course) error {
	return DB.Create(course).Error
}

func GetCourseByID(id uint) (*model.Course, error) {
	var course model.Course
	err := DB.First(&course, id).Error
	return &course, err
}

// GetCoursesByTeacher retrieves all courses taught by a specific teacher.
func GetCoursesByTeacher(teacherID uint) ([]model.Course, error) {
	var courses []model.Course
	err := DB.Where("teacher_id = ?", teacherID).Preload("Class").Preload("TimeSlot").Find(&courses).Error
	return courses, err
}

// GetCoursesByClass retrieves all courses for a specific class.
func GetCoursesByClass(classID uint) ([]model.Course, error) {
	var courses []model.Course
	err := DB.Where("class_id = ?", classID).Preload("Teacher").Preload("TimeSlot").Find(&courses).Error
	return courses, err
}
