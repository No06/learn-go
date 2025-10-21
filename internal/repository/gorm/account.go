package gormrepo

import (
	"context"

	"learn-go/internal/domain"
	"learn-go/internal/repository"

	"gorm.io/gorm"
)

// AccountStore implements repository.AccountRepository using GORM.
type AccountStore struct {
	db *gorm.DB
}

// NewAccountStore creates a new AccountStore.
func NewAccountStore(db *gorm.DB) *AccountStore {
	return &AccountStore{db: db}
}

func (s *AccountStore) Create(ctx context.Context, account *domain.Account) error {
	return s.db.WithContext(ctx).Create(account).Error
}

func (s *AccountStore) FindByIdentifier(ctx context.Context, schoolID, identifier string) (*domain.Account, error) {
	var account domain.Account
	if err := s.db.WithContext(ctx).Where("school_id = ? AND identifier = ?", schoolID, identifier).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *AccountStore) FindByID(ctx context.Context, id string) (*domain.Account, error) {
	var account domain.Account
	if err := s.db.WithContext(ctx).First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *AccountStore) ListByRole(ctx context.Context, schoolID string, role domain.Role, page, size int) ([]domain.Account, int64, error) {
	var (
		accounts []domain.Account
		total    int64
	)

	query := s.db.WithContext(ctx).Where("school_id = ? AND role = ?", schoolID, role)
	if err := query.Model(&domain.Account{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}
	return accounts, total, nil
}

var _ repository.AccountRepository = (*AccountStore)(nil)
