package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	appErr "victa/internal/errors"
)

type JWTService struct {
	secret []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: []byte(secret)}
}

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

func (s *JWTService) ParseToken(tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, appErr.ErrUnexpectedMethod
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return 0, appErr.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, appErr.ErrInvalidToken
	}

	coStr, ok := claims["company_id"].(string)
	if !ok {
		return 0, appErr.ErrClaimMissing
	}

	companyID, err := strconv.ParseInt(coStr, 10, 64)
	if err != nil || companyID <= 0 {
		return 0, appErr.ErrInvalidToken
	}
	return companyID, nil
}
