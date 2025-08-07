package domain

type AppUpdateType string

const (
	AppUpdateName          AppUpdateType = "name"
	AppUpdateSlug          AppUpdateType = "slug"
	AppUpdateAppStoreURL   AppUpdateType = "app_store_url"
	AppUpdatePlayStoreURL  AppUpdateType = "play_store_url"
	AppUpdateRuStoreURL    AppUpdateType = "ru_store_url"
	AppUpdateAppGalleryURL AppUpdateType = "app_gallery_url"
)
