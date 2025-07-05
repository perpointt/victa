package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildCompanyDetail(chatID int64, company *domain.Company, user *domain.User) tgbotapi.MessageConfig {
	text := b.GetCompanyDetailMessage(company)

	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📱 Приложения", fmt.Sprintf("%s?company_id=%v", CallbackListApp, company.ID)),
	))

	err := b.CompanySvc.CheckAdmin(user.ID, company.ID)

	if err == nil {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 Участники", fmt.Sprintf("%s?company_id=%v", CallbackListUser, company.ID)),
		))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧩 Интеграции", fmt.Sprintf("%s?company_id=%v", CallbackCompanyIntegrations, company.ID)),
		))

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			b.BuildDeleteButton(fmt.Sprintf("%s?company_id=%v", CallbackDeleteCompany, company.ID)),
			b.BuildEditButton(fmt.Sprintf("%s?company_id=%v", CallbackUpdateCompany, company.ID)),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildCloseButton(),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	message := b.NewKeyboardMessage(chatID, text, keyboard)
	return message
}
