package domain

import "time"

// User описывает сущность пользователя.
type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`            // Здесь хранится хэш пароля
	Companies []Company `json:"companies,omitempty"` // Список компаний, с которыми связан пользователь
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
