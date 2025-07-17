package domain

import "time"

type CompanySecret struct {
	CompanyID int64      // владелец
	Type      SecretType // классификатор тайны
	Cipher    []byte     // шифротекст (nonce + data)
	CreatedAt time.Time  // метка вставки / последнего апдейта
}
