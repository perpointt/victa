package domain

import "time"

// User описывает сущность пользователя.
type User struct {
	ID        int64     `json:"id"`
	TgID      string    `json:"tg_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
