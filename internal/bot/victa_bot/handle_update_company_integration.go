package victa_bot

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	appErr "victa/internal/errors"
)

func (b *Bot) HandleUpdateCompanyIntegrationCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверные параметры."))
		return
	}
	companyID := params.CompanyID

	ci, err := b.CompanySvc.GetCompanyIntegrationByID(ctx, companyID)
	if err != nil && !errors.Is(err, appErr.ErrIntegrationNotFound) {
		b.SendErrorMessage(chatID, err)
		return
	}

	tmpl, err := b.BuildIntegrationTemplate(ci)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	msgText := fmt.Sprintf(
		"Отправьте обновленный JSON с данными для интеграции:\n\n```json\n%s\n```",
		tmpl,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			b.BuildCancelButton(),
		),
	)

	b.AddChatState(chatID, StateWaitingUpdateCompanyIntegration)
	b.AddPendingCompanyID(chatID, params.CompanyID)

	b.SendPendingMessage(b.NewKeyboardMessage(chatID, msgText, keyboard))
}

func (b *Bot) HandleUpdateCompanyIntegration(ctx context.Context, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	companyID := b.pendingCompanyIDs[chatID]

	company, err := b.CompanySvc.GetByID(ctx, companyID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	_, err = b.CompanySvc.CreateOrUpdateCompanyIntegration(ctx, company.ID, message.Text)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	config, err := b.BuildCompanyIntegrationsDetail(ctx, chatID, company)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.ClearChatState(chatID)

	b.SendMessage(*config)
}
