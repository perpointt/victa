package service

import (
	"context"
	"fmt"

	"victa/internal/domain"
	"victa/internal/repository"
)

type UserService struct {
	usersRepo     repository.UserRepository
	companiesRepo repository.UserCompanyRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	userCompaniesRepo repository.UserCompanyRepository,
) *UserService {
	return &UserService{
		usersRepo:     userRepo,
		companiesRepo: userCompaniesRepo,
	}
}

func (s *UserService) Register(ctx context.Context, tgID int64, name string) (*domain.User, error) {
	u := &domain.User{TgID: fmt.Sprintf("%d", tgID), Name: name}
	return s.usersRepo.Create(ctx, u)
}

func (s *UserService) Update(ctx context.Context, userID int64, name string) (*domain.User, error) {
	u := &domain.User{ID: userID, Name: name}
	return s.usersRepo.Update(ctx, u)
}

func (s *UserService) GetByTgID(ctx context.Context, tgID int64) (*domain.User, error) {
	return s.usersRepo.GetByTgID(ctx, tgID)
}

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

func (s *UserService) DeleteFromCompany(ctx context.Context, userID, companyID int64) error {
	return s.companiesRepo.Delete(ctx, userID, companyID)
}
