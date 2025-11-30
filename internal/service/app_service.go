package service

import (
	"context"
	"strings"

	"victa/internal/domain"
	appErr "victa/internal/errors"
	"victa/internal/repository"
)

type AppService struct {
	appRepo repository.AppRepository
}

func NewAppService(repo repository.AppRepository) *AppService {
	return &AppService{appRepo: repo}
}

func (s *AppService) GetByID(ctx context.Context, appID int64) (*domain.App, error) {
	return s.appRepo.GetByID(ctx, appID)
}

func (s *AppService) GetAllByCompanyID(ctx context.Context, companyID int64) ([]domain.App, error) {
	return s.appRepo.GetAllByCompanyID(ctx, companyID)
}

func (s *AppService) Create(ctx context.Context, companyID int64, name, slug string) (*domain.App, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(slug) == "" {
		return nil, appErr.ErrInvalidInput
	}

	app := &domain.App{
		CompanyID: companyID,
		Name:      name,
		Slug:      slug,
	}
	return s.appRepo.Create(ctx, app)
}

func (s *AppService) Update(ctx context.Context, app *domain.App) (*domain.App, error) {
	return s.appRepo.Update(ctx, app)
}

func (s *AppService) Delete(ctx context.Context, appID int64) error {
	return s.appRepo.Delete(ctx, appID)
}
