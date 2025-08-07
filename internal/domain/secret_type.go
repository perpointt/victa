package domain

// SecretType — тип строковых констант для секретов.
type SecretType string

const (
	SecretGoogleJSON                 SecretType = "google_json"
	SecretAppleP8                    SecretType = "apple_p8"
	SecretAppleIssuerID              SecretType = "apple_issuer_id"
	SecretAppleKeyID                 SecretType = "apple_key_id"
	SecretCodemagicApiKey            SecretType = "codemagic_api_key"
	SecretNotificationBotToken       SecretType = "notification_bot_token"
	SecretErrorsNotificationChatID   SecretType = "errors_notification_chat_id"
	SecretDeployNotificationChatID   SecretType = "deploy_notification_chat_id"
	SecretIssuesNotificationChatID   SecretType = "issues_notification_chat_id"
	SecretVersionsNotificationChatID SecretType = "versions_notification_chat_id"
	SecretReviewsNotificationChatID  SecretType = "reviews_notification_chat_id"
)
