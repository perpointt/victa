package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildUserList(chatID int64, company *domain.Company) (*tgbotapi.MessageConfig, error) {
	users, err := b.UserSvc.GetAllByCompanyID(company.ID)
	if err != nil {
		return nil, err
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range users {
		cbData := fmt.Sprintf("%v:%d", CallbackDetailUser, c.ID)
		userTitle := fmt.Sprintf("%v (ID: %d)", c.Name, c.ID)
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(userTitle, cbData),
			),
		)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Пригласить пользователя", CallbackInviteUser),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildBackButton(fmt.Sprintf("%v:%d", CallbackBackToDetailCompany, company.ID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := b.NewKeyboardMessage(chatID, fmt.Sprintf(
		"*%s* (ID: %d)",
		company.Name,
		company.ID,
	), keyboard)
	return &msg, nil
}
