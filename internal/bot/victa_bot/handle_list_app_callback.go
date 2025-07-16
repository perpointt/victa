package victa_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleListAppsCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID
	tgID := callback.From.ID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	company, err := b.CompanySvc.GetByID(ctx, params.CompanyID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	message, err := b.BuildAppList(ctx, chatID, tgID, company)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.EditMessage(messageID, *message)
}
