package service

import (
	"victa/internal/domain"
	"victa/internal/repository"
)

// UserService содержит логику работы с пользователями
type UserService struct {
	UserRepo          repository.UserRepository
	UserCompaniesRepo repository.UserCompanyRepository
}

// NewUserService создаёт новый сервис пользователей
func NewUserService(userRepo repository.UserRepository, userCompaniesRepo repository.UserCompanyRepository) *UserService {
	return &UserService{UserRepo: userRepo, UserCompaniesRepo: userCompaniesRepo}
}

// Register регистрирует пользователя (создаёт или обновляет запись)
func (s *UserService) Register(tgID string, name string) (*domain.User, error) {
	u := &domain.User{TgID: tgID, Name: name}
	user, err := s.UserRepo.Create(u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByTgID(tgID int64) (*domain.User, error) {
	user, err := s.UserRepo.GetByTgID(tgID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAllDetailByCompanyID(companyID int64) ([]domain.UserDetail, error) {
	users, err := s.UserRepo.GetAllByCompanyID(companyID)
	companies, err := s.UserCompaniesRepo.GetAllByCompanyID(companyID)

	if err != nil {
		return nil, err
	}

	userMap := make(map[int64]domain.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	var details []domain.UserDetail
	for _, uc := range companies {
		if u, ok := userMap[uc.UserID]; ok {
			details = append(details, domain.UserDetail{
				User:    u,
				Company: uc,
			})
		}
	}

	return details, nil
}

// GetByCompanyAndUserID возвращает детальную информацию по одному пользователю в рамках заданной компании.
func (s *UserService) GetByCompanyAndUserID(companyID, userID int64) (*domain.UserDetail, error) {
	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	uc, err := s.UserCompaniesRepo.GetByCompanyAndUserID(companyID, userID)
	if err != nil {
		return nil, err
	}

	detail := &domain.UserDetail{
		User:    *user,
		Company: *uc,
	}
	return detail, nil
}

func (s *UserService) DeleteFromCompany(userID, companyID int64) error {
	return s.UserCompaniesRepo.Delete(userID, companyID)
}
