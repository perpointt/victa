package domain

import "time"

// App описывает приложение, привязанное к компании.
type App struct {
	ID            int64     `json:"id"`
	CompanyID     int64     `json:"company_id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	AppStoreURL   *string   `json:"app_store_url,omitempty"`
	PlayStoreURL  *string   `json:"play_store_url,omitempty"`
	RuStoreURL    *string   `json:"ru_store_url,omitempty"`
	AppGalleryURL *string   `json:"app_gallery_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
