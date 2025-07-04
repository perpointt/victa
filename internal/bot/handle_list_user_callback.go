package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleListUsersCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID
	idPtr, err := b.GetIdFromCallback(callback.Data)
	if err != nil || idPtr == nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}
	companyID := *idPtr

	company, err := b.CompanySvc.GetById(companyID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при поиске компании."))
		return
	}

	message, err := b.BuildUserList(chatID, company)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Ошибка при построении списка пользователей: %v", err)))
		return
	}

	b.EditMessage(messageID, *message)
}
