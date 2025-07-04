package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildUserDetail(chatID int64, user *domain.UserDetail) tgbotapi.MessageConfig {
	text := fmt.Sprintf(
		"*%s* (ID: %d)\n\n%s",
		user.User.Name,
		user.User.ID,
		b.GetRoleTitle(user.Company.RoleID),
	)

	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildCloseButton(CallbackDeleteMessage),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	message := b.NewKeyboardMessage(chatID, text, keyboard)
	return message
}
