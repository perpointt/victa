package domain

import "time"

// PlayVersion — «x.y.z+n»: Semantic = x.y.z, Code = n.
type PlayVersion struct {
	Semantic string // "1.2.3"
	Code     int64  // 45
}

// PlayReleaseInfo — то, что лежит в production‑треке.
type PlayReleaseInfo struct {
	PackageName string
	Version     PlayVersion
	// при желании можно добавить ReleaseNotes []string и т. д.
}

// PlayReview — один отзыв пользователя.
type PlayReview struct {
	ReviewID     string
	AuthorName   string
	Rating       int64
	Text         string
	LastModified time.Time
}
