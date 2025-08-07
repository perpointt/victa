package victa_bot

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"victa/internal/domain"
	appErr "victa/internal/errors"
)

func (b *Bot) BuildCompanyIntegrationsDetail(
	ctx context.Context,
	chatID int64,
	company *domain.Company,
) (*tgbotapi.MessageConfig, error) {
	secrets, err := b.CompanySvc.GetAllSecretsByCompanyID(ctx, company.ID)
	if err != nil && !errors.Is(err, appErr.ErrSecretNotFound) {
		return nil, err
	}
	has := make(map[domain.SecretType]bool, len(secrets))
	for _, s := range secrets {
		has[s.Type] = true
	}

	mark := func(ok bool) string {
		if ok {
			return "🟢"
		}
		return "🔴"
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	// helper: убираем «notification» из SecretType → короче callback_data
	shortType := func(st domain.SecretType) string {
		return strings.ReplaceAll(string(st), "notification", "")
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s Codemagic API Key", mark(has[domain.SecretCodemagicApiKey])),
			fmt.Sprintf("%s?company_id=%d&secret_type=%s",
				CallbackUpdateCompanyIntegrations,
				company.ID,
				domain.SecretCodemagicApiKey),
		),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s Notification Bot Token", mark(has[domain.SecretNotificationBotToken])),
			fmt.Sprintf("%s?company_id=%d&secret_type=%s",
				CallbackUpdateCompanyIntegrations,
				company.ID,
				shortType(domain.SecretNotificationBotToken)),
		),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s Deploy Notification Chat ID", mark(has[domain.SecretDeployNotificationChatID])),
			fmt.Sprintf("%s?company_id=%d&secret_type=%s",
				CallbackUpdateCompanyIntegrations,
				company.ID,
				shortType(domain.SecretDeployNotificationChatID)),
		),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s Issue Notification Chat ID", mark(has[domain.SecretIssuesNotificationChatID])),
			fmt.Sprintf("%s?company_id=%d&secret_type=%s",
				CallbackUpdateCompanyIntegrations,
				company.ID,
				shortType(domain.SecretIssuesNotificationChatID)),
		),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s Error Notification Chat ID", mark(has[domain.SecretErrorsNotificationChatID])),
			fmt.Sprintf("%s?company_id=%d&secret_type=%s",
				CallbackUpdateCompanyIntegrations,
				company.ID,
				shortType(domain.SecretErrorsNotificationChatID)),
		),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s Versions Notification Chat ID", mark(has[domain.SecretVersionsNotificationChatID])),
			fmt.Sprintf("%s?company_id=%d&secret_type=%s",
				CallbackUpdateCompanyIntegrations,
				company.ID,
				shortType(domain.SecretVersionsNotificationChatID)),
		),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s Reviews Notification Chat ID", mark(has[domain.SecretReviewsNotificationChatID])),
			fmt.Sprintf("%s?company_id=%d&secret_type=%s",
				CallbackUpdateCompanyIntegrations,
				company.ID,
				shortType(domain.SecretReviewsNotificationChatID)),
		),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s Apple P8", mark(has[domain.SecretAppleP8])),
			fmt.Sprintf("%s?company_id=%d&secret_type=%s",
				CallbackUpdateCompanyIntegrations,
				company.ID,
				shortType(domain.SecretAppleP8)),
		),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s Apple Key ID", mark(has[domain.SecretAppleKeyID])),
			fmt.Sprintf("%s?company_id=%d&secret_type=%s",
				CallbackUpdateCompanyIntegrations,
				company.ID,
				shortType(domain.SecretAppleKeyID)),
		),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s Apple Issuer ID", mark(has[domain.SecretAppleIssuerID])),
			fmt.Sprintf("%s?company_id=%d&secret_type=%s",
				CallbackUpdateCompanyIntegrations,
				company.ID,
				shortType(domain.SecretAppleIssuerID)),
		),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s Google JSON", mark(has[domain.SecretGoogleJSON])),
			fmt.Sprintf("%s?company_id=%d&secret_type=%s",
				CallbackUpdateCompanyIntegrations,
				company.ID,
				shortType(domain.SecretGoogleJSON)),
		),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔒 Сгенерировать API токен",
			fmt.Sprintf("%v?company_id=%d", CallbackCreateJwtToken, company.ID)),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildBackButton(fmt.Sprintf("%v?company_id=%d", CallbackBackToDetailCompany, company.ID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	text := fmt.Sprintf("💼 *%s | Интеграции* 🧩", company.Name)
	cfg := b.NewKeyboardMessage(chatID, text, keyboard)

	return &cfg, nil
}
