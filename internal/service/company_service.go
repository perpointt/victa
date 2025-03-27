package service

import (
	"victa/internal/domain"
	"victa/internal/repository"
)

// CompanyService описывает бизнес-логику для компаний.
type CompanyService interface {
	CreateCompany(company *domain.Company, userID int64) error
	GetCompanies(userID int64) ([]domain.Company, error)
	GetCompanyByID(id int64, userID int64) (*domain.Company, error)
	UpdateCompany(company *domain.Company) error
	DeleteCompany(id int64) error
}

type companyService struct {
	repo repository.CompanyRepository
}

// NewCompanyService создаёт новый экземпляр сервиса.
func NewCompanyService(repo repository.CompanyRepository) CompanyService {
	return &companyService{repo: repo}
}

func (s *companyService) CreateCompany(company *domain.Company, userID int64) error {
	return s.repo.CreateCompanyWithUser(company, userID)
}

func (s *companyService) GetCompanies(userID int64) ([]domain.Company, error) {
	return s.repo.GetAllWithUser(userID)
}

func (s *companyService) GetCompanyByID(id int64, userID int64) (*domain.Company, error) {
	return s.repo.GetByIdWithUser(id, userID)
}

func (s *companyService) UpdateCompany(company *domain.Company) error {
	return s.repo.Update(company)
}

func (s *companyService) DeleteCompany(id int64) error {
	return s.repo.Delete(id)
}
