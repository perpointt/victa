package service

import (
	"victa/internal/domain"
	"victa/internal/repository"
)

type CompanyService struct {
	repo repository.CompanyRepository
}

// NewCompanyService создаёт новый экземпляр сервиса.
func NewCompanyService(repo repository.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) GetAllByUserId(userID int64) ([]domain.Company, error) {
	return s.repo.GetAllByUserId(userID)
}

func (s *CompanyService) GetById(companyID int64) (*domain.Company, error) {
	return s.repo.GetById(companyID)
}

func (s *CompanyService) Create(name string, userID int64) (*domain.Company, error) {
	return s.repo.Create(domain.Company{Name: name}, userID)
}
