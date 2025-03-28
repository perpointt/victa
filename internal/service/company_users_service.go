package service

import (
	"victa/internal/repository"
)

type CompanyUserService interface {
	IsUserAdminInCompany(userID, companyID int64) (bool, error)
	IsUserInCompany(userID, companyID int64) (bool, error)
	LinkUsersWithCompany(userIDs []int64, companyID int64) error
	UnlinkUsersCompany(userIDs []int64, companyID int64) error
}

type userCompanyService struct {
	repo repository.UserCompanyRepository
}

func NewUserCompanyService(repo repository.UserCompanyRepository) CompanyUserService {
	return &userCompanyService{repo: repo}
}

func (s *userCompanyService) IsUserAdminInCompany(userID, companyID int64) (bool, error) {
	role, err := s.repo.GetUserRole(userID, companyID)
	if err != nil {
		return false, err
	}
	return role == "admin", nil
}

func (s *userCompanyService) IsUserInCompany(userID, companyID int64) (bool, error) {
	role, err := s.repo.GetUserRole(userID, companyID)
	if err != nil {
		return false, err
	}
	return role != "", nil
}

func (s *userCompanyService) LinkUsersWithCompany(userIDs []int64, companyID int64) error {
	if err := s.repo.LinkUserWithCompany(userIDs, companyID); err != nil {
		return err
	}

	return nil
}

func (s *userCompanyService) UnlinkUsersCompany(userIDs []int64, companyID int64) error {
	if err := s.repo.UnlinkUserFromCompany(userIDs, companyID); err != nil {
		return err
	}

	return nil
}
