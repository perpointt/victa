package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ErrUnexpectedMethod — токен подписан не HS256.
var ErrUnexpectedMethod = errors.New("unexpected signing method")

// ErrInvalidToken — подпись битая или формат не распознан.
var ErrInvalidToken = errors.New("invalid token")

// ErrClaimMissing — нет обязательного claim‑а company_id.
var ErrClaimMissing = errors.New("company_id claim missing")

// JWTService генерирует и валидирует JWT для компаний (без exp по‑умолчанию).
type JWTService struct {
	secret []byte
}

// NewJWTService инициализирует сервис с HMAC‑секретом.
func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: []byte(secret)}
}

/*
GenerateToken

Создаёт HS256‑подписанный токен вида:

	{
	  "company_id": "<id>",
	  "iat":        <unix>
	}

Срок действия не задаётся — токен бессрочный.
Строка возвращается с префиксом «Bearer » для удобства.
*/
func (s *JWTService) GenerateToken(companyID int64) (string, error) {
	claims := jwt.MapClaims{
		"company_id": strconv.FormatInt(companyID, 10),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return "Bearer " + signed, nil
}

/*
ParseToken

Проверяет подпись, достаёт company_id и возвращает его.
Ошибки маппятся на публичные Err* из верхней части файла.
*/
func (s *JWTService) ParseToken(tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrUnexpectedMethod
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}

	coStr, ok := claims["company_id"].(string)
	if !ok {
		return 0, ErrClaimMissing
	}

	companyID, err := strconv.ParseInt(coStr, 10, 64)
	if err != nil || companyID <= 0 {
		return 0, ErrInvalidToken
	}
	return companyID, nil
}
