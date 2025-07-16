package repository

import (
	"context"
	"victa/internal/domain"
)

type AppRepository interface {
	GetByID(ctx context.Context, appID int64) (*domain.App, error)
	GetAllByCompanyID(ctx context.Context, companyID int64) ([]domain.App, error)
	Create(ctx context.Context, app *domain.App) (*domain.App, error)
	Update(ctx context.Context, app *domain.App) (*domain.App, error)
	Delete(ctx context.Context, appID int64) error
}
