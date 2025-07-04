package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"time"
)

// InviteService умеет генерировать и валидировать приглашения без БД.
type InviteService struct {
	secret []byte        // секрет для HMAC
	ttl    time.Duration // время жизни токена
}

// NewInviteService создаёт InviteService с TTL 48h.
func NewInviteService(secret []byte) *InviteService {
	return &InviteService{
		secret: secret,
		ttl:    48 * time.Hour,
	}
}

// CreateToken генерирует токен-приглашение для компании.
func (s *InviteService) CreateToken(companyID int64) string {
	// payload: 8 байт companyID + 8 байт unix timestamp
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[0:8], uint64(companyID))
	exp := time.Now().Add(s.ttl).Unix()
	binary.BigEndian.PutUint64(buf[8:16], uint64(exp))

	// HMAC
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(buf)
	sig := mac.Sum(nil)

	// token = base64(payload||sig)
	token := append(buf, sig...)
	return base64.URLEncoding.EncodeToString(token)
}

// ValidateToken проверяет токен и возвращает companyID или ошибку.
func (s *InviteService) ValidateToken(token string) (int64, error) {
	raw, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return 0, errors.New("invalid token encoding")
	}
	if len(raw) != 16+sha256.Size {
		return 0, errors.New("invalid token length")
	}
	payload := raw[:16]
	sig := raw[16:]

	// проверяем HMAC
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(payload)
	expected := mac.Sum(nil)
	if !hmac.Equal(sig, expected) {
		return 0, errors.New("invalid token signature")
	}

	// извлекаем companyID и expire
	companyID := int64(binary.BigEndian.Uint64(payload[0:8]))
	exp := int64(binary.BigEndian.Uint64(payload[8:16]))

	if time.Now().Unix() > exp {
		return 0, errors.New("token expired")
	}
	return companyID, nil
}
