package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildUserDetail(chatID int64, user *domain.UserDetail) tgbotapi.MessageConfig {
	text := fmt.Sprintf("👤 *%s*\n\n%s\n*Роль*: %s",
		user.User.Name,
		b.GetUserDetailMessage(&user.User),
		b.GetRoleTitle(user.Company.RoleID),
	)

	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить из компании", fmt.Sprintf("%v?user_id=%d&company_id=%d", CallbackDeleteUser, user.User.ID, user.Company.CompanyID)),
	))

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildCloseButton(CallbackDeleteMessage),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	message := b.NewKeyboardMessage(chatID, text, keyboard)
	return message
}
