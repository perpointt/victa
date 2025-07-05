package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildCompanyIntegrationsDetail(chatID int64, company *domain.Company) (*tgbotapi.MessageConfig, error) {
	ci, err := b.CompanySvc.GetCompanyIntegrationByID(company.ID)
	if err != nil {
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

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			b.BuildEditButton(fmt.Sprintf("%s?company_id=%d", CallbackUpdateCompanyIntegrations, company.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			b.BuildBackButton(fmt.Sprintf("%v?company_id=%d", CallbackBackToDetailCompany, company.ID)),
		),
	)

	config := b.NewKeyboardMessage(chatID, text, keyboard)

	return &config, nil
}
