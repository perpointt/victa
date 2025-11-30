package repository

import (
	"context"
	"victa/internal/domain"
)

type CompanySecretRepository interface {
	Create(ctx context.Context, secret *domain.CompanySecret) (*domain.CompanySecret, error)
	GetByCompanyIDAndType(ctx context.Context, companyID int64, secretType domain.SecretType) (*domain.CompanySecret, error)
	GetAllByCompanyID(ctx context.Context, companyID int64) ([]domain.CompanySecret, error)
}
