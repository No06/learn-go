package gormrepo

import (
	"context"

	"learn-go/internal/domain"
	"learn-go/internal/repository"

	"gorm.io/gorm"
)

// AssignmentStore implements repository.AssignmentRepository with GORM.
type AssignmentStore struct {
	db *gorm.DB
}

func NewAssignmentStore(db *gorm.DB) *AssignmentStore {
	return &AssignmentStore{db: db}
}

func (s *AssignmentStore) Create(ctx context.Context, assignment *domain.Assignment, questions []domain.AssignmentQuestion) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(assignment).Error; err != nil {
			return err
		}
		for i := range questions {
			questions[i].AssignmentID = assignment.ID
		}
		if len(questions) > 0 {
			if err := tx.Create(&questions).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *AssignmentStore) Get(ctx context.Context, id string) (*domain.Assignment, []domain.AssignmentQuestion, error) {
	var assignment domain.Assignment
	if err := s.db.WithContext(ctx).First(&assignment, "id = ?", id).Error; err != nil {
		return nil, nil, err
	}
	var questions []domain.AssignmentQuestion
	if err := s.db.WithContext(ctx).Where("assignment_id = ?", assignment.ID).Order("order_index").Find(&questions).Error; err != nil {
		return nil, nil, err
	}
	return &assignment, questions, nil
}

var _ repository.AssignmentRepository = (*AssignmentStore)(nil)
