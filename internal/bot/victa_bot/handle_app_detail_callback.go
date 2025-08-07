package victa_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleDetailAppCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	message, err := b.CreateAppDetailMessage(ctx, callback)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.ClearChatState(chatID)
	b.SendMessage(*message)
}

func (b *Bot) HandleBackToDetailAppCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	message, err := b.CreateAppDetailMessage(ctx, callback)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.ClearChatState(chatID)
	b.EditMessage(messageID, *message)
}

func (b *Bot) CreateAppDetailMessage(ctx context.Context, callback *tgbotapi.CallbackQuery) (*tgbotapi.MessageConfig, error) {
	chatID := callback.Message.Chat.ID
	tgID := callback.From.ID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		return nil, err
	}

	app, err := b.AppSvc.GetByID(ctx, params.AppID)
	if err != nil {
		return nil, err
	}

	user, err := b.UserSvc.GetByTgID(ctx, tgID)
	if err != nil {
		return nil, err
	}

	detail := b.BuildAppDetail(ctx, chatID, app, user)
	return &detail, nil
}
