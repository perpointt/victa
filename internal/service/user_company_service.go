package service

import "victa/internal/repository"

type UserCompanyService interface {
	// IsAdmin возвращает true, если пользователь имеет роль "admin" в компании.
	IsAdmin(userID, companyID int64) (bool, error)
}

type userCompanyService struct {
	repo repository.UserCompanyRepository
}

func NewUserCompanyService(repo repository.UserCompanyRepository) UserCompanyService {
	return &userCompanyService{repo: repo}
}

func (s *userCompanyService) IsAdmin(userID, companyID int64) (bool, error) {
	role, err := s.repo.GetUserRole(userID, companyID)
	if err != nil {
		return false, err
	}
	return role == "admin", nil
}
