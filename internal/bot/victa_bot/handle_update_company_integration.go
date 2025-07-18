package victa_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"victa/internal/domain"
)

func (b *Bot) HandleUpdateCompanyIntegrationCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверные параметры."))
		return
	}

	// helper: убираем «notification» из SecretType → короче callback_data
	shortType := func(st domain.SecretType) string {
		return strings.ReplaceAll(string(st), "notification", "")
	}

	switch string(params.SecretType) {
	case string(domain.SecretCodemagicApiKey):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateCodemagicApiKey)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Codemagic API Key", keyboard))
	case shortType(domain.SecretNotificationBotToken):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateNotificationBotToken)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Notification Bot Token", keyboard))
	case shortType(domain.SecretDeployNotificationChatID):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateDeployNotificationChatID)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Deploy Notification Chat ID", keyboard))
	case shortType(domain.SecretIssuesNotificationChatID):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateIssueNotificationChatID)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Issue Notification Chat ID", keyboard))
	case shortType(domain.SecretErrorsNotificationChatID):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateErrorNotificationChatID)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Error Notification Chat ID", keyboard))
	case shortType(domain.SecretAppleP8):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUploadAppleP8)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте AppleP8", keyboard))
	case shortType(domain.SecretGoogleJSON):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUploadGoogleJSON)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Google JSON", keyboard))
	}
}

func (b *Bot) HandleUpdateCompanySecret(ctx context.Context, message *tgbotapi.Message, secretType domain.SecretType) {
	chatID := message.Chat.ID
	companyID := b.pendingCompanyIDs[chatID]

	company, err := b.CompanySvc.GetByID(ctx, companyID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	// helper: убираем «notification» из SecretType → короче callback_data
	shortType := func(st domain.SecretType) string {
		return strings.ReplaceAll(string(st), "notification", "")
	}

	switch string(secretType) {
	case shortType(domain.SecretAppleP8),
		shortType(domain.SecretGoogleJSON):
		// TODO
	default:
		_, err = b.CompanySvc.CreateTextSecret(ctx, company.ID, secretType, message.Text)
		if err != nil {
			b.SendErrorMessage(chatID, err)
			return
		}
	}

	config, err := b.BuildCompanyIntegrationsDetail(ctx, chatID, company)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.DeleteMessage(chatID, message.MessageID)
	b.ClearChatState(chatID)
	b.SendMessage(*config)
}
