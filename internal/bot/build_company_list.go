package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildCompanyList(chatID int64, user *domain.User) (*tgbotapi.MessageConfig, error) {
	companies, err := b.CompanySvc.GetAllByUserId(user.ID)
	if err != nil {
		return nil, err
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range companies {
		cbData := fmt.Sprintf("%v?company_id=%d", CallbackDetailCompany, c.ID)
		companyTitle := fmt.Sprintf("%v (ID: %d)", c.Name, c.ID)
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(companyTitle, cbData),
			),
		)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Создать компанию", CallbackCreateCompany),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := b.NewKeyboardMessage(chatID, "", keyboard)
	return &msg, nil
}
