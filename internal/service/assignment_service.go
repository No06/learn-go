package service

import (
	"time"

	"hinoob.net/learn-go/internal/database"
	"hinoob.net/learn-go/internal/model"
)

// --- Assignment Services ---

// CreateAssignmentForTeacher handles the logic for a teacher creating an assignment.
func CreateAssignmentForTeacher(title, description string, dueDate time.Time, teacherID uint) (*model.Assignment, error) {
	assignment := &model.Assignment{
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		CreatedByID: teacherID,
	}
	err := database.CreateAssignment(assignment)
	return assignment, err
}

// --- Submission Services ---

// CreateOrUpdateSubmission handles the logic for a student submitting their work.
func CreateOrUpdateSubmission(assignmentID, studentID uint, content string, fileURLs []string) (*model.Submission, error) {
	// Create the main submission record
	submission := &model.Submission{
		AssignmentID: assignmentID,
		StudentID:    studentID,
		Content:      content,
	}
	if err := database.CreateSubmission(submission); err != nil {
		return nil, err
	}

	// If file URLs are provided, create the submission file records
	if len(fileURLs) > 0 {
		files := make([]model.SubmissionFile, len(fileURLs))
		for i, url := range fileURLs {
			files[i] = model.SubmissionFile{
				SubmissionID: submission.ID,
				FileURL:      url,
			}
		}
		// This is not the most efficient way (multiple inserts).
		// A better approach for production would be a bulk insert.
		if err := database.CreateSubmissionFiles(files); err != nil {
			// In a real app, you might want to roll back the submission creation here.
			return nil, err
		}
	}

	return database.GetSubmissionByID(submission.ID)
}

// GradeAndCommentOnSubmission handles the logic for a teacher grading a submission.
func GradeAndCommentOnSubmission(submissionID, teacherID uint, grade, commentText string, isImage bool) (*model.Submission, error) {
	// 1. Get the submission
	submission, err := database.GetSubmissionByID(submissionID)
	if err != nil {
		return nil, err
	}

	// 2. Update the grade
	submission.Grade = grade
	if err := database.UpdateSubmission(submission); err != nil {
		return nil, err
	}

	// 3. Add the comment if provided
	if commentText != "" {
		comment := &model.Comment{
			SubmissionID: submissionID,
			UserID:       teacherID,
			Content:      commentText,
			IsImage:      isImage,
		}
		if err := database.CreateComment(comment); err != nil {
			return nil, err
		}
	}

	// 4. Return the updated submission with all comments
	return database.GetSubmissionByID(submissionID)
}
