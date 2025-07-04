package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleDetailUserCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	message, err := b.CreateUserDetailMessage(callback)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при получении данных пользователя."))
		return
	}

	b.ClearChatState(chatID)
	b.SendMessage(*message)
}

func (b *Bot) HandleBackToDetailUserCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	message, err := b.CreateUserDetailMessage(callback)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при получении данных пользователя."))
		return
	}

	b.ClearChatState(chatID)
	b.EditMessage(messageID, *message)
}

func (b *Bot) CreateUserDetailMessage(callback *tgbotapi.CallbackQuery) (*tgbotapi.MessageConfig, error) {
	chatID := callback.Message.Chat.ID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		return nil, err
	}

	detailUser, err := b.UserSvc.GetByCompanyAndUserID(params.CompanyID, params.UserID)
	if err != nil {
		return nil, err
	}

	detail := b.BuildUserDetail(chatID, detailUser)
	return &detail, nil
}
