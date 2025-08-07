package victa_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleCompanyIntegrationsCallback обрабатывает нажатие кнопки «Интеграции»
func (b *Bot) HandleCompanyIntegrationsCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверные параметры."))
		return
	}
	companyID := params.CompanyID

	company, err := b.CompanySvc.GetByID(ctx, companyID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	config, err := b.BuildCompanyIntegrationsDetail(ctx, chatID, company)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}
	b.EditMessage(messageID, *config)
}
