package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTService генерирует и проверяет JWT-токены для компаний (без срока действия).
type JWTService struct {
	secret []byte
}

// NewJWTService создаёт JWTService с заданным секретом.
func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: []byte(secret)}
}

// GenerateToken создаёт JWT для указанного companyID.
// Токен не содержит claim-а exp, поэтому действует бессрочно.
func (s *JWTService) GenerateToken(companyID int64) (string, error) {
	claims := jwt.MapClaims{
		"company_id": strconv.FormatInt(companyID, 10),
		"iat":        time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ParseToken разбирает и проверяет токен, возвращает companyID.
func (s *JWTService) ParseToken(tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}
	coStr, ok := claims["company_id"].(string)
	if !ok {
		return 0, errors.New("company_id claim missing")
	}
	companyID, err := strconv.ParseInt(coStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid company_id claim")
	}
	return companyID, nil
}
