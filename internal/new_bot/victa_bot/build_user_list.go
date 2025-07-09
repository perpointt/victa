package victa_bot

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
			suffix = " (Это вы)"
			cbData = CallbackBlank
		}

		title := fmt.Sprintf("👤 %s | %s%s",
			c.User.Name, b.GetRoleTitle(c.Company.RoleID), suffix,
		)

		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(title, cbData),
			),
		)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ Пригласить пользователя", fmt.Sprintf("%v?company_id=%d", CallbackInviteUser, company.ID)),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildBackButton(fmt.Sprintf("%v?company_id=%d", CallbackBackToDetailCompany, company.ID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	text := fmt.Sprintf("💼 *%s | Участники* 👥", company.Name)

	msg := b.NewKeyboardMessage(chatID, text, keyboard)
	return &msg, nil
}
