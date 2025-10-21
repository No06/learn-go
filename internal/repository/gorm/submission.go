package gormrepo

import (
	"context"

	"learn-go/internal/domain"
	"learn-go/internal/repository"

	"gorm.io/gorm"
)

// SubmissionStore implements repository.SubmissionRepository with GORM.
type SubmissionStore struct {
	db *gorm.DB
}

func NewSubmissionStore(db *gorm.DB) *SubmissionStore {
	return &SubmissionStore{db: db}
}

func (s *SubmissionStore) CreateOrUpdate(ctx context.Context, submission *domain.AssignmentSubmission, items []domain.SubmissionItem) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if submission.ID == "" {
			if err := tx.Create(submission).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Save(submission).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("submission_id = ?", submission.ID).Delete(&domain.SubmissionItem{}).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].SubmissionID = submission.ID
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *SubmissionStore) ListByAssignment(ctx context.Context, assignmentID string) ([]domain.AssignmentSubmission, error) {
	var submissions []domain.AssignmentSubmission
	if err := s.db.WithContext(ctx).Where("assignment_id = ?", assignmentID).Order("updated_at DESC").Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (s *SubmissionStore) ListItemsBySubmissionIDs(ctx context.Context, submissionIDs []string) ([]domain.SubmissionItem, error) {
	if len(submissionIDs) == 0 {
		return nil, nil
	}
	var items []domain.SubmissionItem
	if err := s.db.WithContext(ctx).Where("submission_id IN ?", submissionIDs).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *SubmissionStore) GetByAssignmentAndStudent(ctx context.Context, assignmentID, studentID string) (*domain.AssignmentSubmission, []domain.SubmissionItem, error) {
	var submission domain.AssignmentSubmission
	if err := s.db.WithContext(ctx).
		Where("assignment_id = ? AND student_id = ?", assignmentID, studentID).
		First(&submission).Error; err != nil {
		return nil, nil, err
	}
	var items []domain.SubmissionItem
	if err := s.db.WithContext(ctx).Where("submission_id = ?", submission.ID).Find(&items).Error; err != nil {
		return nil, nil, err
	}
	return &submission, items, nil
}

func (s *SubmissionStore) GetByID(ctx context.Context, submissionID string) (*domain.AssignmentSubmission, []domain.SubmissionItem, error) {
	var submission domain.AssignmentSubmission
	if err := s.db.WithContext(ctx).First(&submission, "id = ?", submissionID).Error; err != nil {
		return nil, nil, err
	}
	var items []domain.SubmissionItem
	if err := s.db.WithContext(ctx).Where("submission_id = ?", submission.ID).Find(&items).Error; err != nil {
		return nil, nil, err
	}
	return &submission, items, nil
}

func (s *SubmissionStore) UpdateGrades(ctx context.Context, submission *domain.AssignmentSubmission, items []domain.SubmissionItem) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.AssignmentSubmission{}).
			Where("id = ?", submission.ID).
			Updates(map[string]interface{}{
				"score":        submission.Score,
				"feedback":     submission.Feedback,
				"status":       submission.Status,
				"updated_at":   submission.UpdatedAt,
				"submitted_at": submission.SubmittedAt,
			}).Error; err != nil {
			return err
		}
		for _, item := range items {
			if err := tx.Model(&domain.SubmissionItem{}).
				Where("id = ? AND submission_id = ?", item.ID, submission.ID).
				Updates(map[string]interface{}{
					"score":  item.Score,
					"answer": item.Answer,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

var _ repository.SubmissionRepository = (*SubmissionStore)(nil)
