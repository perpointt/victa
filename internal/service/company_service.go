package service

import (
	"context"
	"fmt"
	"victa/internal/crypto"
	"victa/internal/domain"
	appErr "victa/internal/errors"
	"victa/internal/repository"
)

// ID роли «admin» в таблице roles. Нулевой элемент справочника.
const adminRoleID int64 = 1

// CompanyService инкапсулирует бизнес‑логику для сущности Company.
type CompanyService struct {
	companyRepo       repository.CompanyRepository
	companySecretRepo repository.CompanySecretRepository
	encryptor         crypto.Encryptor
}

// NewCompanyService создаёт экземпляр сервиса компаний.
func NewCompanyService(
	companyRepo repository.CompanyRepository,
	companySecretRepo repository.CompanySecretRepository,
	encryptor crypto.Encryptor,
) *CompanyService {
	return &CompanyService{
		companyRepo:       companyRepo,
		companySecretRepo: companySecretRepo,
		encryptor:         encryptor,
	}
}

// GetAllByUserID возвращает список компаний, к которым привязан пользователь.
func (s *CompanyService) GetAllByUserID(ctx context.Context, userID int64) ([]domain.Company, error) {
	return s.companyRepo.GetAllByUserID(ctx, userID)
}

// GetByID возвращает одну компанию по ID (или ErrCompanyNotFound из appRepo).
func (s *CompanyService) GetByID(ctx context.Context, companyID int64) (*domain.Company, error) {
	return s.companyRepo.GetByID(ctx, companyID)
}

// Create создаёт новую компанию и сразу делает userID её администратором.
func (s *CompanyService) Create(ctx context.Context, name string, adminUserID int64) (*domain.Company, error) {
	return s.companyRepo.Create(ctx, domain.Company{Name: name}, adminUserID)
}

// Update изменяет название компании. Разрешено только админу.
func (s *CompanyService) Update(ctx context.Context, companyID int64, name string, userID int64) (*domain.Company, error) {
	if err := s.CheckAdmin(ctx, userID, companyID); err != nil {
		return nil, err
	}
	return s.companyRepo.Update(ctx, domain.Company{ID: companyID, Name: name})
}

// Delete удаляет компанию целиком (каскадно), если вызван админом.
func (s *CompanyService) Delete(ctx context.Context, companyID, userID int64) error {
	if err := s.CheckAdmin(ctx, userID, companyID); err != nil {
		return err
	}
	return s.companyRepo.Delete(ctx, companyID)
}

// AddUserToCompany добавляет пользователя в компанию с ролью developer.
func (s *CompanyService) AddUserToCompany(ctx context.Context, userID, companyID int64) error {
	return s.companyRepo.AddUserToCompany(ctx, userID, companyID, "developer")
}

// GetSecret возвращает нужный секрет компании.
func (s *CompanyService) GetSecret(ctx context.Context, companyID int64, secretType domain.SecretType) ([]byte, error) {
	sec, err := s.companySecretRepo.GetByCompanyIDAndType(ctx, companyID, secretType)
	if err != nil {
		return nil, err
	}
	plain, err := s.encryptor.Open(ctx, sec.Cipher)
	if err != nil {
		return nil, err
	}
	return plain, nil
}

// GetAllSecretsByCompanyID возвращает нужный секрет компании.
func (s *CompanyService) GetAllSecretsByCompanyID(ctx context.Context, companyID int64) ([]domain.CompanySecret, error) {
	secrets, err := s.companySecretRepo.GetAllByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

// CreateTextSecret шифрует произвольный текст и сохраняет как секрет компании.
func (s *CompanyService) CreateTextSecret(
	ctx context.Context,
	companyID int64,
	secretType domain.SecretType,
	text string,
) (*domain.CompanySecret, error) {
	cipher, err := s.encryptor.Seal(ctx, []byte(text))
	if err != nil {
		return nil, fmt.Errorf("encrypt secret: %w", err)
	}

	sec := &domain.CompanySecret{
		CompanyID: companyID,
		Type:      secretType,
		Cipher:    cipher,
	}

	created, err := s.companySecretRepo.Create(ctx, sec)
	if err != nil {
		return nil, fmt.Errorf("store secret: %w", err)
	}
	return created, nil
}

func (s *CompanyService) CreateBinarySecret(
	ctx context.Context,
	companyID int64,
	secretType domain.SecretType,
	data []byte,
) (*domain.CompanySecret, error) {

	if len(data) == 0 {
		return nil, appErr.ErrInvalidInput
	}
	cipher, err := s.encryptor.Seal(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}
	sec := &domain.CompanySecret{
		CompanyID: companyID,
		Type:      secretType,
		Cipher:    cipher,
	}
	return s.companySecretRepo.Create(ctx, sec)
}

// CheckAdmin проверяет, что userID имеет роль admin в заданной компании.
// Возвращает ErrNotCompanyAdmin, если прав не хватает.
func (s *CompanyService) CheckAdmin(ctx context.Context, userID, companyID int64) error {
	roleIDPtr, err := s.companyRepo.GetUserRole(ctx, userID, companyID)
	if err != nil {
		return err
	}
	if roleIDPtr == nil || *roleIDPtr != adminRoleID {
		return appErr.ErrNotCompanyAdmin
	}
	return nil
}
