package database

import "hinoob.net/learn-go/internal/model"

// CreateUser creates a new user in the database
func CreateUser(user *model.User) error {
	result := DB.Create(user)
	return result.Error
}

// GetUserByUsername finds a user by their username
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	result := DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByID finds a user by their ID
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	result := DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// AssignTeachersToStudent links one or more teachers to a student
func AssignTeachersToStudent(student *model.User, teachers []*model.User) error {
	return DB.Model(student).Association("Teachers").Append(teachers)
}
