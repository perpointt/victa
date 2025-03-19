package service

import (
	"victa/internal/domain"
	"victa/internal/repository"
)

// AppService описывает бизнес-логику для приложений.
type AppService interface {
	CreateApp(app *domain.App) error
	GetAllApps() ([]domain.App, error)
	GetAppByID(id int64) (*domain.App, error)
	UpdateApp(app *domain.App) error
	DeleteApp(id int64) error
}

type appService struct {
	repo repository.AppRepository
}

// NewAppService создаёт новый экземпляр AppService.
func NewAppService(repo repository.AppRepository) AppService {
	return &appService{repo: repo}
}

func (s *appService) CreateApp(app *domain.App) error {
	return s.repo.Create(app)
}

func (s *appService) GetAllApps() ([]domain.App, error) {
	return s.repo.GetAll()
}

func (s *appService) GetAppByID(id int64) (*domain.App, error) {
	return s.repo.GetByID(id)
}

func (s *appService) UpdateApp(app *domain.App) error {
	return s.repo.Update(app)
}

func (s *appService) DeleteApp(id int64) error {
	return s.repo.Delete(id)
}
