package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
	appErr "victa/internal/errors"
)

const (
	expiryLen = 6
	sigLen    = 6
)

var now = time.Now

type InviteService struct {
	secret []byte
	ttl    time.Duration
}

func NewInviteService(secret []byte, ttl time.Duration) *InviteService {
	return &InviteService{secret: secret, ttl: ttl}
}

func (s *InviteService) CreateToken(companyID int64) (string, error) {
	if companyID <= 0 {
		return "", appErr.ErrBadFormat
	}

	cmp36 := strconv.FormatInt(companyID, 36)

	exp36 := strconv.FormatInt(now().Add(s.ttl).Unix(), 36)
	exp36 = leftPad(exp36, expiryLen, '0')

	msg := cmp36 + exp36
	sig := firstNHex(hmacSHA256(s.secret, msg), sigLen)

	return msg + sig, nil
}

func (s *InviteService) ValidateToken(token string) (int64, error) {
	if len(token) < expiryLen+sigLen+1 { // хотя бы 1 символ на companyID
		return 0, appErr.ErrBadFormat
	}

	msgEnd := len(token) - sigLen
	expStart := msgEnd - expiryLen

	sig := token[msgEnd:]
	exp36 := token[expStart:msgEnd]
	cmp36 := token[:expStart]

	expUnix, err := strconv.ParseInt(exp36, 36, 64)
	if err != nil {
		return 0, appErr.ErrBadFormat
	}
	if now().Unix() > expUnix {
		return 0, appErr.ErrExpired
	}

	expected := firstNHex(hmacSHA256(s.secret, cmp36+exp36), sigLen)
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return 0, appErr.ErrBadSignature
	}

	companyID, err := strconv.ParseInt(cmp36, 36, 64)
	if err != nil || companyID <= 0 {
		return 0, appErr.ErrBadFormat
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
