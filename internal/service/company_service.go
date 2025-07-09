package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"victa/internal/domain"
	"victa/internal/repository"
)

type CompanyService struct {
	CompanyRepo     repository.CompanyRepository
	IntegrationRepo repository.CompanyIntegrationRepository
}

// NewCompanyService создаёт новый экземпляр сервиса.
func NewCompanyService(
	companyRepository repository.CompanyRepository,
	integrationRepo repository.CompanyIntegrationRepository,
) *CompanyService {
	return &CompanyService{
		CompanyRepo:     companyRepository,
		IntegrationRepo: integrationRepo,
	}
}

func (s *CompanyService) GetAllByUserID(userID int64) ([]domain.Company, error) {
	return s.CompanyRepo.GetAllByUserID(userID)
}

func (s *CompanyService) GetByID(companyID int64) (*domain.Company, error) {
	return s.CompanyRepo.GetByID(companyID)
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

// GetCompanyIntegrationByID возвращает настройки интеграций для компании.
func (s *CompanyService) GetCompanyIntegrationByID(companyID int64) (*domain.CompanyIntegration, error) {
	return s.IntegrationRepo.GetByID(companyID)
}

// CreateOrUpdateCompanyIntegration парсит JSON-пayload и сохраняет интеграции.
// Требует, чтобы userID был администратором компании.
func (s *CompanyService) CreateOrUpdateCompanyIntegration(
	companyID int64,
	payload string,
) (*domain.CompanyIntegration, error) {
	var ci domain.CompanyIntegration
	if err := json.Unmarshal([]byte(payload), &ci); err != nil {
		return nil, fmt.Errorf("неверный формат JSON интеграций: %w", err)
	}

	ci.CompanyID = companyID

	return s.IntegrationRepo.CreateOrUpdate(&ci)
}
