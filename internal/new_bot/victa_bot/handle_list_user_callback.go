package victa_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleListUsersCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID
	tgID := callback.From.ID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	company, err := b.CompanySvc.GetByID(params.CompanyID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при поиске компании."))
		return
	}

	message, err := b.BuildUserList(chatID, tgID, company)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Ошибка при построении списка пользователей: %v", err)))
		return
	}

	b.EditMessage(messageID, *message)
}
