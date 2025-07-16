package service

import (
	"context"
	"errors"
	"strings"

	"victa/internal/domain"
	"victa/internal/repository"
)

// ErrInvalidInput возвращается при пустом name или slug.
var ErrInvalidInput = errors.New("app name and slug must be non-empty")

// AppService инкапсулирует бизнес‑логику для сущности App.
type AppService struct {
	repo repository.AppRepository
}

// NewAppService создаёт новый сервис для работы с приложениями.
func NewAppService(repo repository.AppRepository) *AppService {
	return &AppService{repo: repo}
}

// GetByID возвращает приложение по ID.
func (s *AppService) GetByID(ctx context.Context, appID int64) (*domain.App, error) {
	return s.repo.GetByID(ctx, appID)
}

// GetAllByCompanyID возвращает все приложения компании.
func (s *AppService) GetAllByCompanyID(ctx context.Context, companyID int64) ([]domain.App, error) {
	return s.repo.GetAllByCompanyID(ctx, companyID)
}

// Create валидирует входные данные и создаёт приложение.
func (s *AppService) Create(ctx context.Context, companyID int64, name, slug string) (*domain.App, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(slug) == "" {
		return nil, ErrInvalidInput
	}

	app := &domain.App{
		CompanyID: companyID,
		Name:      name,
		Slug:      slug,
	}
	return s.repo.Create(ctx, app)
}

// Update изменяет имя и slug приложения.
func (s *AppService) Update(ctx context.Context, appID int64, name, slug string) (*domain.App, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(slug) == "" {
		return nil, ErrInvalidInput
	}

	app := &domain.App{
		ID:   appID,
		Name: name,
		Slug: slug,
	}
	return s.repo.Update(ctx, app)
}

// Delete удаляет приложение по ID.
func (s *AppService) Delete(ctx context.Context, appID int64) error {
	return s.repo.Delete(ctx, appID)
}
