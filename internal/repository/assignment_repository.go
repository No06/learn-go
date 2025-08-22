package repository

import "hinoob.net/learn-go/internal/model"

// --- Assignment Repository ---

// CreateAssignment creates a new assignment in the database.
func CreateAssignment(assignment *model.Assignment) error {
	return DB.Create(assignment).Error
}

// GetAssignmentByID retrieves an assignment by its ID.
func GetAssignmentByID(id uint) (*model.Assignment, error) {
	var assignment model.Assignment
	err := DB.Preload("CreatedBy").First(&assignment, id).Error
	return &assignment, err
}

// GetAssignmentsByTeacherID retrieves all assignments created by a specific teacher.
func GetAssignmentsByTeacherID(teacherID uint) ([]model.Assignment, error) {
	var assignments []model.Assignment
	err := DB.Where("created_by_id = ?", teacherID).Find(&assignments).Error
	return assignments, err
}

// --- Submission Repository ---

// CreateSubmission creates a new submission in the database.
func CreateSubmission(submission *model.Submission) error {
	return DB.Create(submission).Error
}

// CreateSubmissionFiles bulk-creates submission file records.
func CreateSubmissionFiles(files []model.SubmissionFile) error {
	return DB.Create(&files).Error
}

// GetSubmissionByID retrieves a submission by its ID, including related data.
func GetSubmissionByID(id uint) (*model.Submission, error) {
	var submission model.Submission
	err := DB.Preload("Student").Preload("Comments.User").Preload("Files").First(&submission, id).Error
	return &submission, err
}

// GetSubmissionsByAssignmentID retrieves all submissions for a given assignment.
func GetSubmissionsByAssignmentID(assignmentID uint) ([]model.Submission, error) {
	var submissions []model.Submission
	err := DB.Where("assignment_id = ?", assignmentID).Preload("Student").Find(&submissions).Error
	return submissions, err
}

// UpdateSubmission updates an existing submission (e.g., for grading).
func UpdateSubmission(submission *model.Submission) error {
	return DB.Save(submission).Error
}

// CreateComment adds a new comment to a submission.
func CreateComment(comment *model.Comment) error {
	return DB.Create(comment).Error
}
