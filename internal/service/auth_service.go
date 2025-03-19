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
	Register(email, password string, companyID *int64) (*domain.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo        repository.UserRepository
	companyRepo     repository.CompanyRepository
	userCompanyRepo repository.UserCompanyRepository
	jwtSecret       string
}

// NewAuthService создаёт новый экземпляр AuthService.
func NewAuthService(userRepo repository.UserRepository, companyRepo repository.CompanyRepository, userCompanyRepo repository.UserCompanyRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:        userRepo,
		companyRepo:     companyRepo,
		userCompanyRepo: userCompanyRepo,
		jwtSecret:       jwtSecret,
	}
}

// Register регистрирует нового пользователя.
// Если companyID указан, то проверяется наличие соответствующего агентства в БД и создаётся связь между пользователем и агентством.
func (s *authService) Register(email, password string, companyID *int64) (*domain.User, error) {
	// Проверяем, существует ли уже пользователь с таким email.
	if user, err := s.userRepo.GetByEmail(email); err == nil && user != nil {
		return nil, errors.New("user already exists")
	}

	// Хэшируем пароль.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаём пользователя без привязки к компании.
	newUser := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	// Если указан companyID, проверяем, существует ли агентство, и связываем его с пользователем.
	if companyID != nil {
		_, err := s.companyRepo.GetByID(*companyID)
		if err != nil {
			return nil, errors.New("company not found")
		}
		if err := s.userCompanyRepo.LinkUserCompany(newUser.ID, *companyID); err != nil {
			return nil, err
		}
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

	// В JWT-токене теперь хранится только user_id и срок действия.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
