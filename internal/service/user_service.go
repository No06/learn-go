package service

import (
	"errors"

	"hinoob.net/learn-go/internal/model"
	"hinoob.net/learn-go/internal/pkg/utils"
	"hinoob.net/learn-go/internal/repository"

	"gorm.io/gorm"
)

// CreateUser handles the business logic for creating a user
func CreateUser(username, password, fullName string, role model.Role, teacherIDs []uint) (*model.User, error) {
	// 1. Check if user already exists
	_, err := repository.GetUserByUsername(username)
	if err == nil {
		return nil, errors.New("username already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// A different error occurred
		return nil, err
	}

	// 2. Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 3. Create the user object
	user := &model.User{
		Username:     username,
		PasswordHash: hashedPassword,
		FullName:     fullName,
		Role:         role,
	}

	// 4. Save the user to the database
	if err := repository.CreateUser(user); err != nil {
		return nil, err
	}

	// 5. If the user is a student and teacher IDs are provided, create the association
	if user.IsStudent() && len(teacherIDs) > 0 {
		var teachers []*model.User
		for _, teacherID := range teacherIDs {
			teacher, err := repository.GetUserByID(teacherID)
			if err != nil {
				// In a real app, you might want to handle this more gracefully
				// (e.g., log a warning, or fail the whole transaction)
				continue
			}
			if teacher.IsTeacher() {
				teachers = append(teachers, teacher)
			}
		}
		if len(teachers) > 0 {
			if err := repository.AssignTeachersToStudent(user, teachers); err != nil {
				return nil, err
			}
		}
	}

	return user, nil
}

// Login handles the user login logic
func Login(username, password string) ( /* token */ string, *model.User, error) {
	// 1. Find user by username
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("user not found")
		}
		return "", nil, err
	}

	// 2. Check password
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, errors.New("invalid password")
	}

	// 3. Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, user, nil
}
