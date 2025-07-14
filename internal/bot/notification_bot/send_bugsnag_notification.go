package notification_bot

func (bot *Bot) SendBugsnagNotification(payload string) {
	text := payload
	_ = bot.SendMessage(bot.NewHtmlMessage(bot.chatID, text))
}
