package domain

import "time"

// User описывает сущность пользователя.
type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Хэш пароля не будет включен в JSON
	Companies []Company `json:"companies,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
