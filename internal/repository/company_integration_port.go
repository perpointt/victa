package repository

import (
	"context"
	"victa/internal/domain"
)

type CompanyIntegrationRepository interface {
	GetByID(ctx context.Context, companyID int64) (*domain.CompanyIntegration, error)
	CreateOrUpdate(ctx context.Context, ci *domain.CompanyIntegration) (*domain.CompanyIntegration, error)
}
