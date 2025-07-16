package service

import (
	"context"
	"fmt"

	"victa/internal/domain"
	"victa/internal/repository"
)

// UserService содержит логику работы с пользователями.
type UserService struct {
	usersRepo     repository.UserRepository
	companiesRepo repository.UserCompanyRepository
}

// NewUserService создаёт новый сервис пользователей.
func NewUserService(
	userRepo repository.UserRepository,
	userCompaniesRepo repository.UserCompanyRepository,
) *UserService {
	return &UserService{
		usersRepo:     userRepo,
		companiesRepo: userCompaniesRepo,
	}
}

// Register создаёт нового пользователя или возвращает существующего.
func (s *UserService) Register(ctx context.Context, tgID int64, name string) (*domain.User, error) {
	u := &domain.User{TgID: fmt.Sprintf("%d", tgID), Name: name}
	return s.usersRepo.Create(ctx, u)
}

// Update меняет имя пользователя.
func (s *UserService) Update(ctx context.Context, userID int64, name string) (*domain.User, error) {
	u := &domain.User{ID: userID, Name: name}
	return s.usersRepo.Update(ctx, u)
}

// GetByTgID ищет пользователя по Telegram‑ID.
func (s *UserService) GetByTgID(ctx context.Context, tgID int64) (*domain.User, error) {
	return s.usersRepo.GetByTgID(ctx, tgID)
}

// GetAllDetailByCompanyID возвращает сотрудников компании + их роли.
func (s *UserService) GetAllDetailByCompanyID(ctx context.Context, companyID int64) ([]domain.UserDetail, error) {
	users, err := s.usersRepo.GetAllByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	relations, err := s.companiesRepo.GetAllByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	userMap := make(map[int64]domain.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	details := make([]domain.UserDetail, 0, len(relations))
	for _, rel := range relations {
		if u, ok := userMap[rel.UserID]; ok {
			details = append(details, domain.UserDetail{
				User:    u,
				Company: rel,
			})
		}
	}
	return details, nil
}

// GetByCompanyAndUserID возвращает детальную инфу по сотруднику в компании.
func (s *UserService) GetByCompanyAndUserID(ctx context.Context, companyID, userID int64) (*domain.UserDetail, error) {
	user, err := s.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	rel, err := s.companiesRepo.GetByCompanyAndUserID(ctx, companyID, userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserDetail{
		User:    *user,
		Company: *rel,
	}, nil
}

// DeleteFromCompany удаляет сотрудника из компании.
func (s *UserService) DeleteFromCompany(ctx context.Context, userID, companyID int64) error {
	return s.companiesRepo.Delete(ctx, userID, companyID)
}
