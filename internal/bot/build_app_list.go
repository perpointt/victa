package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildAppList(chatID, tgID int64, company *domain.Company) (*tgbotapi.MessageConfig, error) {
	apps, err := b.AppSvc.GetAllByCompanyID(company.ID)
	if err != nil {
		return nil, err
	}

	user, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		return nil, err
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range apps {
		cbData := fmt.Sprintf("%v?app_id=%d", CallbackDetailApp, c.ID)
		title := fmt.Sprintf("📱 %s | %s", c.Name, c.Slug)
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(title, cbData),
			),
		)
	}

	err = b.CompanySvc.CheckAdmin(user.ID, company.ID)
	if err == nil {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать приложение", fmt.Sprintf("%v?company_id=%d", CallbackCreateApp, company.ID)),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildBackButton(fmt.Sprintf("%v?company_id=%d", CallbackBackToDetailCompany, company.ID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	text := fmt.Sprintf("💼 *%s | Приложения* 📱", company.Name)

	msg := b.NewKeyboardMessage(chatID, text, keyboard)
	return &msg, nil
}
