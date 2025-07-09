package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"
)

// InviteService генерирует короткие токены без точек:
// формат "<companyID base36><expiry base36 6 chars><sig 6 hex chars>"
type InviteService struct {
	secret []byte
	ttl    time.Duration
}

// NewInviteService создаёт InviteService с TTL 48h.
func NewInviteService(secret []byte) *InviteService {
	return &InviteService{secret: secret, ttl: 48 * time.Hour}
}

// CreateToken возвращает токен в виде непрерывной строки:
//   - первые N символов — companyID в base36,
//   - следующие 6 символов — expiry UNIX в base36, слева нули,
//   - последние 6 символов — первые 6 hex байт HMAC(prefix).
func (s *InviteService) CreateToken(companyID int64) string {
	// companyID в base36
	cmp := strconv.FormatInt(companyID, 36)

	// expiry в base36, ровно 6 символов
	exp := strconv.FormatInt(time.Now().Add(s.ttl).Unix(), 36)
	if len(exp) < 6 {
		exp = strings.Repeat("0", 6-len(exp)) + exp
	}

	// подпись HMAC-SHA256(prefix) первые 6 байт → 12 hex, но берем 6 hex
	msg := cmp + exp
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(msg))
	sig := hex.EncodeToString(mac.Sum(nil))[:6]

	return msg + sig
}

// ValidateToken извлекает companyID или ошибку из токена без точек.
func (s *InviteService) ValidateToken(token string) (int64, error) {
	if len(token) < 13 { // минимум 1+6+6
		return 0, errors.New("invalid token length")
	}
	// разбиваем: sig = последние 6, exp = предшествующие 6, cmp = всё остальное
	sig := token[len(token)-6:]
	exp36 := token[len(token)-12 : len(token)-6]
	cmp36 := token[:len(token)-12]

	// проверяем expiry
	expUnix, err := strconv.ParseInt(exp36, 36, 64)
	if err != nil {
		return 0, errors.New("invalid expiry")
	}
	if time.Now().Unix() > expUnix {
		return 0, errors.New("token expired")
	}

	// проверяем подпись
	msg := cmp36 + exp36
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(msg))
	expected := hex.EncodeToString(mac.Sum(nil))[:6]
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return 0, errors.New("invalid signature")
	}

	// парсим companyID
	companyID, err := strconv.ParseInt(cmp36, 36, 64)
	if err != nil {
		return 0, errors.New("invalid companyID")
	}
	return companyID, nil
}
