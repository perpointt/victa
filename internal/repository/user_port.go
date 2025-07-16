package repository

import (
	"context"
	"victa/internal/domain"
)

// UserRepository описывает операции хранения пользователей.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByTgID(ctx context.Context, tgID int64) (*domain.User, error)
	GetAllByCompanyID(ctx context.Context, companyID int64) ([]domain.User, error)
}
