package service

import (
	"context"
	"fmt"
	"victa/internal/crypto"
	"victa/internal/domain"
	appErr "victa/internal/errors"
	"victa/internal/repository"
)

const adminRoleID int64 = 1

type CompanyService struct {
	companyRepo       repository.CompanyRepository
	companySecretRepo repository.CompanySecretRepository
	encryptor         crypto.Encryptor
}

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

func (s *CompanyService) GetAllByUserID(ctx context.Context, userID int64) ([]domain.Company, error) {
	return s.companyRepo.GetAllByUserID(ctx, userID)
}

func (s *CompanyService) GetByID(ctx context.Context, companyID int64) (*domain.Company, error) {
	return s.companyRepo.GetByID(ctx, companyID)
}

func (s *CompanyService) Create(ctx context.Context, name string, adminUserID int64) (*domain.Company, error) {
	return s.companyRepo.Create(ctx, domain.Company{Name: name}, adminUserID)
}

func (s *CompanyService) Update(ctx context.Context, companyID int64, name string, userID int64) (*domain.Company, error) {
	if err := s.CheckAdmin(ctx, userID, companyID); err != nil {
		return nil, err
	}
	return s.companyRepo.Update(ctx, domain.Company{ID: companyID, Name: name})
}

func (s *CompanyService) Delete(ctx context.Context, companyID, userID int64) error {
	if err := s.CheckAdmin(ctx, userID, companyID); err != nil {
		return err
	}
	return s.companyRepo.Delete(ctx, companyID)
}

func (s *CompanyService) AddUserToCompany(ctx context.Context, userID, companyID int64) error {
	return s.companyRepo.AddUserToCompany(ctx, userID, companyID, "developer")
}

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

func (s *CompanyService) GetAllSecretsByCompanyID(ctx context.Context, companyID int64) ([]domain.CompanySecret, error) {
	secrets, err := s.companySecretRepo.GetAllByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

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
func (s *CompanyService) GetAll(ctx context.Context) ([]domain.Company, error) {
	return s.companyRepo.GetAll(ctx)
}
