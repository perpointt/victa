package service

import (
	"errors"
	"victa/internal/domain"
	"victa/internal/repository"
)

type CompanyService struct {
	CompanyRepo repository.CompanyRepository
}

// NewCompanyService создаёт новый экземпляр сервиса.
func NewCompanyService(companyRepository repository.CompanyRepository) *CompanyService {
	return &CompanyService{CompanyRepo: companyRepository}
}

func (s *CompanyService) GetAllByUserId(userID int64) ([]domain.Company, error) {
	return s.CompanyRepo.GetAllByUserId(userID)
}

func (s *CompanyService) GetById(companyID int64) (*domain.Company, error) {
	return s.CompanyRepo.GetById(companyID)
}

func (s *CompanyService) Create(name string, userID int64) (*domain.Company, error) {
	return s.CompanyRepo.Create(domain.Company{Name: name}, userID)
}

// Update изменяет название компании, только если userID — админ в ней.
func (s *CompanyService) Update(companyID int64, name string, userID int64) (*domain.Company, error) {
	if err := s.CheckAdmin(userID, companyID); err != nil {
		return nil, err
	}
	return s.CompanyRepo.Update(domain.Company{ID: companyID, Name: name})
}

func (s *CompanyService) Delete(companyID, userID int64) error {
	if err := s.CheckAdmin(userID, companyID); err != nil {
		return err
	}
	return s.CompanyRepo.Delete(companyID)
}

func (s *CompanyService) CheckAdmin(userID, companyID int64) error {
	roleIDPtr, err := s.CompanyRepo.GetUserRole(userID, companyID)
	if err != nil {
		return err
	}
	// Если связи нет или роль не равна 1 (admin)
	if roleIDPtr == nil || *roleIDPtr != 1 {
		return errors.New("операция доступна только администратору компании")
	}
	return nil
}

func (s *CompanyService) AddUserToCompany(userID, companyID int64) error {
	return s.CompanyRepo.AddUserToCompany(userID, companyID, "developer")
}
