package domain

type AppStoreReleaseInfo struct {
	AppID    string      // 10‑значный App ID
	BundleID string      // com.example.app
	Version  PlayVersion // Semantic + Code (build number)
}

type AppStoreReview = PlayReview // структура из play.go подходит 1‑в‑1
