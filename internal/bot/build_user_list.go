package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"victa/internal/domain"
)

func (b *Bot) BuildUserList(chatID, tgID int64, company *domain.Company) (*tgbotapi.MessageConfig, error) {
	users, err := b.UserSvc.GetAllDetailByCompanyID(company.ID)
	if err != nil {
		return nil, err
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range users {
		userTgID, _ := strconv.ParseInt(c.User.TgID, 10, 64)
		cbData := fmt.Sprintf("%v?user_id=%d&company_id=%d", CallbackDetailUser, c.User.ID, c.Company.CompanyID)
		suffix := ""
		if userTgID == tgID {
			suffix = " (Вы)"
			cbData = CallbackBlank
		}

		userTitle := fmt.Sprintf("%s (ID: %d) | %s%s",
			c.User.Name, c.User.ID, b.GetRoleTitle(c.Company.RoleID), suffix,
		)

		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(userTitle, cbData),
			),
		)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Пригласить пользователя", fmt.Sprintf("%v?company_id=%d", CallbackInviteUser, company.ID)),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildBackButton(fmt.Sprintf("%v?company_id=%d", CallbackBackToDetailCompany, company.ID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := b.NewKeyboardMessage(chatID, fmt.Sprintf(
		"*%s* (ID: %d)",
		company.Name,
		company.ID,
	), keyboard)
	return &msg, nil
}
