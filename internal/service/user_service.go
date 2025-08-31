package service

import (
	"errors"
	"gorm.io/gorm"
	"hinoob.net/learn-go/internal/model"
	"hinoob.net/learn-go/internal/pkg/hash"
	"hinoob.net/learn-go/internal/pkg/jwt"
	"hinoob.net/learn-go/internal/repository"
)

// CreateUser handles the business logic for creating a user
func CreateUser(username, pass, fullName string, role model.Role, teacherIDs []uint) (*model.User, error) {
	// 1. Check if user already exists
	_, err := repository.GetUserByUsername(username)
	if err == nil {
		return nil, errors.New("username already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// A different error occurred
		return nil, err
	}

	// 2. Hash the hash
	hashedPassword, err := hash.HashPassword(pass)
	if err != nil {
		return nil, errors.New("failed to hash hash")
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
func Login(username, pass string) ( /* token */ string, *model.User, error) {
	// 1. Find user by username
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("user not found")
		}
		return "", nil, err
	}

	// 2. Check hash
	if !hash.CheckPasswordHash(pass, user.PasswordHash) {
		return "", nil, errors.New("invalid hash")
	}

	// 3. Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username, string(user.Role))
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, user, nil
}
