package gormrepo

import (
	"context"

	"learn-go/internal/domain"
	"learn-go/internal/repository"

	"gorm.io/gorm"
)

// ClassStore implements repository.ClassRepository using GORM.
type ClassStore struct {
	db *gorm.DB
}

func NewClassStore(db *gorm.DB) *ClassStore {
	return &ClassStore{db: db}
}

func (s *ClassStore) Create(ctx context.Context, class *domain.Class) error {
	return s.db.WithContext(ctx).Create(class).Error
}

func (s *ClassStore) ListByDepartment(ctx context.Context, schoolID, departmentID string) ([]domain.Class, error) {
	var classes []domain.Class
	query := s.db.WithContext(ctx).Where("school_id = ?", schoolID)
	if departmentID != "" {
		query = query.Where("department_id = ?", departmentID)
	}
	if err := query.Order("created_at").Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

func (s *ClassStore) GetByID(ctx context.Context, id string) (*domain.Class, error) {
	var class domain.Class
	if err := s.db.WithContext(ctx).First(&class, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

var _ repository.ClassRepository = (*ClassStore)(nil)
