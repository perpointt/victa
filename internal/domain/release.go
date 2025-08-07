package domain

import "time"

// StoreType — идентификатор «откуда» пришли данные.
type StoreType string

const (
	StoreGooglePlay StoreType = "google_play"
	StoreAppStore   StoreType = "app_store"
)

// ReleaseInfo описывает один production-релиз.
type ReleaseInfo struct {
	Store      StoreType // из какого стора
	AppID      string    // packageName (Play) или AppID (App Store)
	BundleID   string    // только для App Store, Play оставляем пустым
	Semantic   string    // "1.2.3"
	Code       int64     // versionCode или build number
	ReleasedAt time.Time // опционально (если появится в API)
}

// Review описывает один отзыв.
type Review struct {
	Store        StoreType
	AppID        string
	ReviewID     string
	AuthorName   string
	Rating       int // 1–5
	Text         string
	LastModified time.Time
}
