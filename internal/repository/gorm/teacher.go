package gormrepo

import (
	"context"

	"learn-go/internal/domain"
	"learn-go/internal/repository"

	"gorm.io/gorm"
)

// TeacherStore implements repository.TeacherRepository using GORM.
type TeacherStore struct {
	db *gorm.DB
}

// NewTeacherStore returns a new teacher store.
func NewTeacherStore(db *gorm.DB) *TeacherStore {
	return &TeacherStore{db: db}
}

func (s *TeacherStore) Create(ctx context.Context, teacher *domain.Teacher) error {
	return s.db.WithContext(ctx).Create(teacher).Error
}

func (s *TeacherStore) GetByNumber(ctx context.Context, schoolID, number string) (*domain.Teacher, error) {
	var teacher domain.Teacher
	if err := s.db.WithContext(ctx).Where("school_id = ? AND number = ?", schoolID, number).First(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (s *TeacherStore) GetByID(ctx context.Context, id string) (*domain.Teacher, error) {
	var teacher domain.Teacher
	if err := s.db.WithContext(ctx).First(&teacher, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

var _ repository.TeacherRepository = (*TeacherStore)(nil)
