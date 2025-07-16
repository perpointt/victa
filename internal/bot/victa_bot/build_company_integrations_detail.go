package victa_bot

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
	appErr "victa/internal/errors"
)

func (b *Bot) BuildCompanyIntegrationsDetail(ctx context.Context, chatID int64, company *domain.Company) (*tgbotapi.MessageConfig, error) {
	ci, err := b.CompanySvc.GetCompanyIntegrationByID(ctx, company.ID)
	if err != nil && !errors.Is(err, appErr.ErrIntegrationNotFound) {
		return nil, err
	}

	tmpl, err := b.BuildIntegrationTemplate(ci)
	if err != nil {
		return nil, err
	}

	var text = fmt.Sprintf("💼 *%s | Интеграции* 🧩", company.Name)
	if ci == nil {
		text = fmt.Sprintf("%s\n\n%s", text, "🔴 Интеграции не настроены")
	} else {
		text = fmt.Sprintf("%s\n\n%s\n\n```json\n%s\n```", text, "🟢 Интеграции настроены", tmpl)
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildEditButton(fmt.Sprintf("%s?company_id=%d", CallbackUpdateCompanyIntegrations, company.ID)),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔒 Сгенерировать API токен", fmt.Sprintf("%v?company_id=%d", CallbackCreateJwtToken, company.ID)),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildBackButton(fmt.Sprintf("%v?company_id=%d", CallbackBackToDetailCompany, company.ID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	config := b.NewKeyboardMessage(chatID, text, keyboard)

	return &config, nil
}
