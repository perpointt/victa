package domain

import "time"

// AppMarketState — последнее известное состояние приложения в конкретном маркете.
type AppMarketState struct {
	AppID        int64   // apps.id
	MarketID     int64   // app_markets.id
	LastVersion  *int64  // последний известный build/versionCode
	LastReviewID *string // последний известный review_id (строка)
	UpdatedAt    time.Time
}
