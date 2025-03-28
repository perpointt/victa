package service

import (
	"victa/internal/domain"
	"victa/internal/repository"
)

// UserService описывает бизнес-логику для пользователей.
type UserService interface {
	GetUsersInCompany(companyID int64) ([]domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService создаёт новый экземпляр UserService.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUsersInCompany(companyID int64) ([]domain.User, error) {
	return s.repo.GetUsersByCompanyID(companyID)
}
