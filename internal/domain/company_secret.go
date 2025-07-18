package domain

import "time"

type CompanySecret struct {
	CompanyID int64      `json:"company_id"`
	Type      SecretType `json:"type"`
	Cipher    []byte     `json:"cipher"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
