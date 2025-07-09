package domain

type CompanyIntegration struct {
	CompanyID            int64   `json:"company_id"`
	CodemagicAPIKey      *string `json:"codemagic_api_key"`
	NotificationBotToken *string `json:"notification_bot_token"`
	NotificationChatID   *string `json:"notification_chat_id"`
}
