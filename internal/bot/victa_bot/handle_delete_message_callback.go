package victa_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleDeleteMessageCallback(cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID

	b.ClearChatState(cb.Message.Chat.ID)
	b.DeleteMessage(chatID, messageID)
}
