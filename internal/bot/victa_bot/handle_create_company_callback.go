package victa_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleCreateCompanyCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	b.AddChatState(chatID, StateWaitingCreateCompanyName)

	msgText := "Отправьте название компании"
	cancelButton := b.BuildCancelButton()
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(cancelButton))

	b.SendPendingMessage(b.NewKeyboardMessage(chatID, msgText, keyboard))
}

func (b *Bot) HandleCompanyNameCreated(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	tgID := message.From.ID

	user, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	company, err := b.CompanySvc.Create(message.Text, user.ID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	config := b.BuildCompanyDetail(chatID, company, user)

	b.ClearChatState(chatID)

	b.SendMessage(config)
}
