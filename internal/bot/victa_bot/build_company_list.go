package victa_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildCompanyList(chatID int64, user *domain.User) (*tgbotapi.MessageConfig, error) {
	companies, err := b.CompanySvc.GetAllByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range companies {
		cbData := fmt.Sprintf("%v?company_id=%d", CallbackDetailCompany, c.ID)
		title := fmt.Sprintf("💼 %v", c.Name)
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(title, cbData),
			),
		)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ Создать компанию", CallbackCreateCompany),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := b.NewKeyboardMessage(chatID, "", keyboard)
	return &msg, nil
}
