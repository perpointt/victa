package domain

type CompanyIntegration struct {
	CompanyID                int64   `json:"company_id"`
	CodemagicAPIKey          *string `json:"codemagic_api_key"`
	NotificationBotToken     *string `json:"notification_bot_token"`
	DeployNotificationChatID *string `json:"deploy_notification_chat_id"`
	IssuesNotificationChatID *string `json:"issues_notification_chat_id"`
	ErrorsNotificationChatID *string `json:"errors_notification_chat_id"`
}
