package service

import (
	"errors"
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

// Update изменяет название компании, только если userID — админ в ней.
func (s *CompanyService) Update(companyID int64, name string, userID int64) (*domain.Company, error) {
	if err := s.checkAdmin(userID, companyID); err != nil {
		return nil, err
	}
	return s.repo.Update(domain.Company{ID: companyID, Name: name})
}

func (s *CompanyService) Delete(companyID, userID int64) error {
	if err := s.checkAdmin(userID, companyID); err != nil {
		return err
	}
	return s.repo.Delete(companyID)
}

func (s *CompanyService) checkAdmin(userID, companyID int64) error {
	slug, err := s.repo.GetUserRole(userID, companyID)
	if err != nil {
		return err
	}
	if slug != "admin" {
		return errors.New("операция доступна только администратору компании")
	}
	return nil
}
