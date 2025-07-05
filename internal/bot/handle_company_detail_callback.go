package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleDetailCompanyCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	message, err := b.CreateCompanyDetailMessage(callback)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при получении данных компании."))
		return
	}

	b.ClearChatState(chatID)
	b.SendMessage(*message)
}

func (b *Bot) HandleBackToDetailCompanyCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	message, err := b.CreateCompanyDetailMessage(callback)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при получении данных компании."))
		return
	}

	b.ClearChatState(chatID)
	b.EditMessage(messageID, *message)
}

func (b *Bot) CreateCompanyDetailMessage(callback *tgbotapi.CallbackQuery) (*tgbotapi.MessageConfig, error) {
	chatID := callback.Message.Chat.ID
	tgID := callback.From.ID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		return nil, err
	}

	company, err := b.CompanySvc.GetByID(params.CompanyID)
	if err != nil {
		return nil, err
	}

	user, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		return nil, err
	}

	detail := b.BuildCompanyDetail(chatID, company, user)
	return &detail, nil
}
