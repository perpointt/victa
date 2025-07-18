package repository

import (
	"context"
	"victa/internal/domain"
)

type CompanySecretRepository interface {
	// Create сохраняет шифротекст; если запись для (companyID, secretType) уже есть —
	// перезаписывает cipher (idempotent).
	Create(ctx context.Context, secret *domain.CompanySecret) (*domain.CompanySecret, error)

	// GetByCompanyIDAndType возвращает шифротекст; если записи нет — appErr.ErrNotFound.
	GetByCompanyIDAndType(ctx context.Context, companyID int64, secretType domain.SecretType) (*domain.CompanySecret, error)

	// GetAllByCompanyID возвращает все секреты компании.
	GetAllByCompanyID(ctx context.Context, companyID int64) ([]domain.CompanySecret, error)
}
