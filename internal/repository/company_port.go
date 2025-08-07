package repository

import (
	"context"
	"victa/internal/domain"
)

type CompanyRepository interface {
	Create(ctx context.Context, company domain.Company, adminUserID int64) (*domain.Company, error)
	Update(ctx context.Context, company domain.Company) (*domain.Company, error)
	Delete(ctx context.Context, companyID int64) error

	GetAllByUserID(ctx context.Context, userID int64) ([]domain.Company, error)
	GetByID(ctx context.Context, companyID int64) (*domain.Company, error)

	GetUserRole(ctx context.Context, userID, companyID int64) (*int64, error)
	AddUserToCompany(ctx context.Context, userID, companyID int64, roleSlug string) error

	GetAll(ctx context.Context) ([]domain.Company, error)
}
