package service

import (
	"victa/internal/domain"
	"victa/internal/repository"
)

type AppService struct {
	repo repository.AppRepository
}

// NewAppService создаёт новый сервис для работы с приложениями.
func NewAppService(repo repository.AppRepository) *AppService {
	return &AppService{repo: repo}
}

// GetByID возвращает приложение по его ID или nil, если не найдено.
func (s *AppService) GetByID(id int64) (*domain.App, error) {
	return s.repo.GetByID(id)
}

// GetAllByCompanyID возвращает все приложения для заданной компании.
func (s *AppService) GetAllByCompanyID(companyID int64) ([]domain.App, error) {
	return s.repo.GetAllByCompanyID(companyID)
}

// Create создаёт новое приложение и возвращает его сущность.
func (s *AppService) Create(companyID int64, name, slug string) (*domain.App, error) {
	app := &domain.App{
		CompanyID: companyID,
		Name:      name,
		Slug:      slug,
	}
	return s.repo.Create(app)
}

// Update изменяет имя и slug приложения и возвращает обновлённую сущность.
func (s *AppService) Update(id int64, name, slug string) (*domain.App, error) {
	app := &domain.App{
		ID:   id,
		Name: name,
		Slug: slug,
	}
	return s.repo.Update(app)
}

// Delete удаляет приложение по его ID.
func (s *AppService) Delete(id int64) error {
	return s.repo.Delete(id)
}
