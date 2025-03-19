package service

import (
	"errors"
	"time"

	"victa/internal/domain"
	"victa/internal/repository"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthService описывает методы аутентификации.
type AuthService interface {
	Register(email, password string, companyID int64) (*domain.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewAuthService создаёт новый экземпляр AuthService.
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register регистрирует нового пользователя, хэшируя пароль перед сохранением.
func (s *authService) Register(email, password string, companyID int64) (*domain.User, error) {
	// Проверяем, существует ли уже пользователь с таким email
	if user, err := s.userRepo.GetByEmail(email); err == nil && user != nil {
		return nil, errors.New("user already exists")
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		CompanyID: companyID,
		Email:     email,
		Password:  string(hashedPassword),
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// Login выполняет аутентификацию и возвращает JWT, если email и пароль верны.
func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Сравниваем хэшированный пароль с предоставленным паролем
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Создаем токен с полезными данными (claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"company_id": user.CompanyID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
