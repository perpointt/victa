package victa_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleMainMenuCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	user, err := b.UserSvc.GetByTgID(callback.From.ID)

	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при проверке пользователя."))
		return
	}

	if user == nil {
		b.SendMessage(b.NewMessage(chatID, "Сначала зарегистрируйтесь через /start."))
		return
	}

	config := b.BuildMainMenu(chatID, user)
	if config == nil {
		return
	}

	b.ClearChatState(chatID)
	b.EditMessage(messageID, *config)
}
