package service

import (
	"victa/internal/domain"
	"victa/internal/repository"
)

// UserService содержит логику работы с пользователями
type UserService struct {
	UserRepo repository.UserRepository
}

// NewUserService создаёт новый сервис пользователей
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{UserRepo: repo}
}

// Register регистрирует пользователя (создаёт или обновляет запись)
func (s *UserService) Register(tgId string, name string) (*domain.User, error) {
	u := &domain.User{TgId: tgId, Name: name}
	user, err := s.UserRepo.Create(u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindByTgID(tgID int64) (*domain.User, error) {
	user, err := s.UserRepo.FindByTgID(tgID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
