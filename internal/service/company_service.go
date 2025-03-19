package service

import (
	"victa/internal/domain"
	"victa/internal/repository"
)

// CompanyService описывает бизнес-логику для компаний.
type CompanyService interface {
	CreateCompany(company *domain.Company) error
	GetAllCompanies() ([]domain.Company, error)
	GetCompanyByID(id int64) (*domain.Company, error)
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

func (s *companyService) CreateCompany(company *domain.Company) error {
	return s.repo.Create(company)
}

func (s *companyService) GetAllCompanies() ([]domain.Company, error) {
	return s.repo.GetAll()
}

func (s *companyService) GetCompanyByID(id int64) (*domain.Company, error) {
	return s.repo.GetByID(id)
}

func (s *companyService) UpdateCompany(company *domain.Company) error {
	return s.repo.Update(company)
}

func (s *companyService) DeleteCompany(id int64) error {
	return s.repo.Delete(id)
}
