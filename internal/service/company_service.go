package service

import (
	"context"
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

// GetCompanySecret возвращает нужный секрет компании.
func (s *CompanyService) GetCompanySecret(ctx context.Context, companyID int64, secretType domain.SecretType) ([]byte, error) {
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

//// CreateOrUpdateCompanyIntegration принимает JSON‑payload,
//// валидирует его и выполняет upsert настроек интеграций.
//func (s *CompanyService) CreateOrUpdateCompanyIntegration(
//	ctx context.Context,
//	companyID int64,
//	payload string,
//) (*domain.CompanyIntegration, error) {
//	var ci domain.CompanyIntegration
//	if err := json.Unmarshal([]byte(payload), &ci); err != nil {
//		return nil, fmt.Errorf("invalid JSON: %w", err)
//	}
//	ci.CompanyID = companyID
//	return s.integrationRepo.CreateOrUpdate(ctx, &ci)
//}

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
