package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildCompanyList(chatID int64, user *domain.User) *tgbotapi.MessageConfig {
	companies, err := b.CompanySvc.GetAllByUserId(user.ID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при получении списка компаний."))
		return nil
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range companies {
		cbData := fmt.Sprintf("%v:%d", CallbackDetailCompany, c.ID)
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(c.Name, cbData),
			),
		)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ Создать компанию", CallbackCreateCompany),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := b.NewKeyboardMessage(chatID, "", keyboard)
	return &msg
}
