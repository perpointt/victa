package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildAppDetail(chatID int64, app *domain.App, user *domain.User) tgbotapi.MessageConfig {
	text := b.GetAppDetailMessage(app)

	var rows [][]tgbotapi.InlineKeyboardButton

	//rows = append(rows, tgbotapi.NewInlineKeyboardRow(
	//	tgbotapi.NewInlineKeyboardButtonData("📱 Приложения", fmt.Sprintf("%s?company_id=%v", CallbackListApp, app.ID)),
	//))

	err := b.CompanySvc.CheckAdmin(user.ID, app.CompanyID)

	if err == nil {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			b.BuildDeleteButton(fmt.Sprintf("%v?app_id=%d&company_id=%d", CallbackDeleteApp, app.ID, app.CompanyID)),
			b.BuildEditButton(fmt.Sprintf("%v?app_id=%d&company_id=%d", CallbackUpdateApp, app.ID, app.CompanyID)),
		))
		//rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData("🧩 Интеграции", fmt.Sprintf("%s?company_id=%v", CallbackCompanyIntegrations, company.ID)),
		//))
		//
		//rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		//	b.BuildDeleteButton(fmt.Sprintf("%s?company_id=%v", CallbackDeleteCompany, company.ID)),
		//	b.BuildEditButton(fmt.Sprintf("%s?company_id=%v", CallbackUpdateCompany, company.ID)),
		//))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildCloseButton(CallbackDeleteMessage),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	message := b.NewKeyboardMessage(chatID, text, keyboard)
	return message
}
