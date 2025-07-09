package victa_bot

import (
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
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		return
	}

	message, err := b.BuildUserList(chatID, tgID, company)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, err.Error()))
		return
	}

	b.EditMessage(messageID, *message)
}
