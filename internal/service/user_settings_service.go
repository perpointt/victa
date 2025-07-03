package service

import (
	"victa/internal/domain"
	"victa/internal/repository"
)

// UserSettingsService содержит логику работы с пользователями
type UserSettingsService struct {
	UserSettingsRepo repository.UserSettingsRepository
}

// NewUserSettingsService создаёт новый сервис пользователей
func NewUserSettingsService(repo repository.UserSettingsRepository) *UserSettingsService {
	return &UserSettingsService{UserSettingsRepo: repo}
}

//func (s *UserSettingsService) Update(tgID int64) (*domain.UserSettings, error) {
//
//}

func (s *UserSettingsService) FindByUserID(userId int64) (*domain.UserSettings, error) {
	return s.UserSettingsRepo.FindByUserID(userId)
}
