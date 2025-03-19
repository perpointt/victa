package domain

import "time"

// User описывает сущность пользователя.
type User struct {
	ID        int64     `json:"id"`
	CompanyID int64     `json:"company_id"` // ID компании, к которой принадлежит пользователь
	Email     string    `json:"email"`
	Password  string    `json:"password"` // Здесь должен храниться хэш пароля
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
