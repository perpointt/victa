package service

import (
	"context"
	"strings"

	"victa/internal/domain"
	appErr "victa/internal/errors"
	"victa/internal/repository"
)

// AppService инкапсулирует бизнес‑логику для сущности App.
type AppService struct {
	appRepo repository.AppRepository
}

// NewAppService создаёт новый сервис для работы с приложениями.
func NewAppService(repo repository.AppRepository) *AppService {
	return &AppService{appRepo: repo}
}

// GetByID возвращает приложение по ID.
func (s *AppService) GetByID(ctx context.Context, appID int64) (*domain.App, error) {
	return s.appRepo.GetByID(ctx, appID)
}

// GetAllByCompanyID возвращает все приложения компании.
func (s *AppService) GetAllByCompanyID(ctx context.Context, companyID int64) ([]domain.App, error) {
	return s.appRepo.GetAllByCompanyID(ctx, companyID)
}

// Create валидирует входные данные и создаёт приложение.
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

// Update изменяет имя и slug приложения.
func (s *AppService) Update(ctx context.Context, appID int64, name, slug string) (*domain.App, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(slug) == "" {
		return nil, appErr.ErrInvalidInput
	}

	app := &domain.App{
		ID:   appID,
		Name: name,
		Slug: slug,
	}
	return s.appRepo.Update(ctx, app)
}

// Delete удаляет приложение по ID.
func (s *AppService) Delete(ctx context.Context, appID int64) error {
	return s.appRepo.Delete(ctx, appID)
}
