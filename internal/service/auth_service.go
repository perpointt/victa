package service

import (
	"errors"
	"fmt"
	"time"

	"victa/internal/domain"
	"victa/internal/repository"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthService описывает методы аутентификации.
type AuthService interface {
	Register(email, password string, companyID *int64) (*domain.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo    repository.UserRepository
	companyRepo repository.CompanyRepository
	jwtSecret   string
}

// NewAuthService создаёт новый экземпляр AuthService.
func NewAuthService(userRepo repository.UserRepository, companyRepo repository.CompanyRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:    userRepo,
		companyRepo: companyRepo,
		jwtSecret:   jwtSecret,
	}
}

// Register регистрирует нового пользователя.
// Если companyID не указан, создаётся новая компания и её ID используется для нового пользователя.
func (s *authService) Register(email, password string, companyID *int64) (*domain.User, error) {
	// Проверяем, существует ли уже пользователь с таким email
	if user, err := s.userRepo.GetByEmail(email); err == nil && user != nil {
		return nil, errors.New("user already exists")
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Если companyID не указан, создаём новую компанию
	var finalCompanyID int64
	if companyID == nil {
		newCompany := &domain.Company{
			Name: fmt.Sprintf("%s's Company", email),
		}
		if err := s.companyRepo.Create(newCompany); err != nil {
			return nil, err
		}
		finalCompanyID = newCompany.ID
	} else {
		finalCompanyID = *companyID
	}

	newUser := &domain.User{
		CompanyID: &finalCompanyID, // Используем указатель на finalCompanyID
		Email:     email,
		Password:  string(hashedPassword),
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// Login выполняет аутентификацию и возвращает JWT-токен при успешной проверке.
func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

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
