package victa_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleDetailAppCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	message, err := b.CreateAppDetailMessage(callback)
	if err != nil {
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		return
	}

	b.ClearChatState(chatID)
	b.SendMessage(*message)
}

func (b *Bot) HandleBackToDetailAppCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	message, err := b.CreateAppDetailMessage(callback)
	if err != nil {
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		return
	}

	b.ClearChatState(chatID)
	b.EditMessage(messageID, *message)
}

func (b *Bot) CreateAppDetailMessage(callback *tgbotapi.CallbackQuery) (*tgbotapi.MessageConfig, error) {
	chatID := callback.Message.Chat.ID
	tgID := callback.From.ID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		return nil, err
	}

	company, err := b.AppSvc.GetByID(params.AppID)
	if err != nil {
		return nil, err
	}

	user, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		return nil, err
	}

	detail := b.BuildAppDetail(chatID, company, user)
	return &detail, nil
}
