package domain

import "time"

// App описывает сущность приложения.
type App struct {
	ID        int64     `json:"id"`
	CompanyID int64     `json:"company_id"` // ID компании, к которой принадлежит приложение
	Name      string    `json:"name"`
	Platform  string    `json:"platform"`  // Например: ios, android и т.д.
	StoreURL  string    `json:"store_url"` // URL в магазине (может быть пустым)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
