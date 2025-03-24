package service

import (
	"victa/internal/domain"
	"victa/internal/repository"
)

// UserService описывает бизнес-логику для пользователей.
type UserService interface {
	//CreateUser(user *domain.User) error
	GetAllUsers() ([]domain.User, error)
	GetAllUsersByCompany(companyID int64) ([]domain.User, error)
	GetUserByID(id int64) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id int64) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService создаёт новый экземпляр UserService.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

//func (s *userService) CreateUser(user *domain.User) error {
//	return s.repo.CreateWithCompany(user, 0)
//}

func (s *userService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAll()
}

func (s *userService) GetAllUsersByCompany(companyID int64) ([]domain.User, error) {
	return s.repo.GetAllByCompanyID(companyID)
}

func (s *userService) GetUserByID(id int64) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) UpdateUser(user *domain.User) error {
	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id int64) error {
	return s.repo.Delete(id)
}
