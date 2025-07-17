package crypto

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"os"

	"golang.org/x/crypto/nacl/secretbox"
	appErr "victa/internal/errors"
)

// Encryptor шифрует и расшифровывает данные для хранения в БД.
// Реализует алгоритм NaCl secretbox: nonce(24) + ciphertext.
type Encryptor struct {
	key [32]byte
}

// NewEncryptor читает master‑ключ из файла (32‑байт hex/raw).
func NewEncryptor(masterKeyPath string) (*Encryptor, error) {
	raw, err := os.ReadFile(masterKeyPath)
	if err != nil {
		return nil, err
	}
	if len(raw) != 32 {
		return nil, errors.New("master key must be 32 bytes")
	}
	var k [32]byte
	copy(k[:], raw)
	return &Encryptor{key: k}, nil
}

// Seal шифрует plain → nonce||cipher; контекст зарезервирован «на будущее».
func (e *Encryptor) Seal(_ context.Context, plain []byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}
	return secretbox.Seal(nonce[:], plain, &nonce, &e.key), nil
}

// Open расшифровывает cipher → plain.
func (e *Encryptor) Open(_ context.Context, cipher []byte) ([]byte, error) {
	if len(cipher) < 24 {
		return nil, appErr.ErrDecrypt
	}
	var nonce [24]byte
	copy(nonce[:], cipher[:24])
	plain, ok := secretbox.Open(nil, cipher[24:], &nonce, &e.key)
	if !ok {
		return nil, appErr.ErrDecrypt
	}
	return plain, nil
}
