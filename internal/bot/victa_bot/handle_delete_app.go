package victa_bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleDeleteAppCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	b.AddChatState(chatID, StateWaitingConfirmDeleteApp)

	msgText := "Подтвердите удаление приложения"
	confirmMessage := b.BuildConfirmMessage(chatID, msgText, fmt.Sprintf("%s?app_id=%v&company_id=%v", CallbackConfirmOperation, params.AppID, params.CompanyID))

	b.SendPendingMessage(confirmMessage)
}

func (b *Bot) HandleConfirmDeleteAppCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
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

	if err := b.AppSvc.Delete(ctx, params.AppID); err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	config, err := b.BuildAppList(ctx, chatID, tgID, company)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.SendMessage(*config)
}
