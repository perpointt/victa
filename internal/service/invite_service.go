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

// ErrBadFormat — токен повреждён: длина, base‑36, либо companyID ≤ 0.
var ErrBadFormat = errors.New("token format is invalid")

// ErrExpired — токен просрочен по полю expiry.
var ErrExpired = errors.New("token expired")

// ErrBadSignature — HMAC‑подпись не совпала с ожидаемой.
var ErrBadSignature = errors.New("token signature mismatch")

const (
	expiryLen = 6 // 6 символов base‑36 = «YYMMDD» до 2058 г.
	sigLen    = 6 // 6 hex‑символов = 24 бит HMAC‑защиты
)

// now переопределяется в юнит‑тестах → детерминированные результаты.
var now = time.Now

// InviteService генерирует и валидирует короткие инвайт‑токены
// формата <companyID base36><expiry base36 6><sig 6 hex>.
type InviteService struct {
	secret []byte        // HMAC‑ключ
	ttl    time.Duration // срок жизни токена
}

// NewInviteService создаёт сервис с заданным TTL (обычно 48 h).
func NewInviteService(secret []byte, ttl time.Duration) *InviteService {
	return &InviteService{secret: secret, ttl: ttl}
}

/*
CreateToken

Возвращает строку‑токен для companyID:
 1. companyID → base‑36 без лидирующих нулей.
 2. expiry    → (now + ttl) в base‑36, слева паддинг до 6 симв.
 3. sig       → первые 6 hex‑симв. HMAC‑SHA256(companyID+expiry).
*/
func (s *InviteService) CreateToken(companyID int64) (string, error) {
	if companyID <= 0 {
		return "", ErrBadFormat
	}

	cmp36 := strconv.FormatInt(companyID, 36)

	exp36 := strconv.FormatInt(now().Add(s.ttl).Unix(), 36)
	exp36 = leftPad(exp36, expiryLen, '0')

	msg := cmp36 + exp36
	sig := firstNHex(hmacSHA256(s.secret, msg), sigLen)

	return msg + sig, nil
}

/*
ValidateToken

Пошагово проверяет токен и возвращает companyID:
 1. Формат и минимальную длину.
 2. Читаем expiry → не истёк ли?
 3. Сверяем HMAC‑подпись.
 4. Парсим companyID.
*/
func (s *InviteService) ValidateToken(token string) (int64, error) {
	if len(token) < expiryLen+sigLen+1 { // хотя бы 1 символ на companyID
		return 0, ErrBadFormat
	}

	// Разбиваем на cmp / exp / sig.
	msgEnd := len(token) - sigLen
	expStart := msgEnd - expiryLen

	sig := token[msgEnd:]
	exp36 := token[expStart:msgEnd]
	cmp36 := token[:expStart]

	// Expiry.
	expUnix, err := strconv.ParseInt(exp36, 36, 64)
	if err != nil {
		return 0, ErrBadFormat
	}
	if now().Unix() > expUnix {
		return 0, ErrExpired
	}

	// Signature.
	expected := firstNHex(hmacSHA256(s.secret, cmp36+exp36), sigLen)
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return 0, ErrBadSignature
	}

	// companyID.
	companyID, err := strconv.ParseInt(cmp36, 36, 64)
	if err != nil || companyID <= 0 {
		return 0, ErrBadFormat
	}
	return companyID, nil
}

func hmacSHA256(key []byte, msg string) []byte {
	m := hmac.New(sha256.New, key)
	m.Write([]byte(msg))
	return m.Sum(nil)
}

func firstNHex(data []byte, n int) string {
	return hex.EncodeToString(data)[:n]
}

func leftPad(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(string(pad), length-len(s)) + s
}
