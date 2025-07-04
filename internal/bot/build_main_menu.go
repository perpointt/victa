package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildMainMenu(chatID int64, user *domain.User) *tgbotapi.MessageConfig {
	text := fmt.Sprintf(
		"Привет, *%s*!\n\n*Список компаний:*",
		user.Name,
	)

	msg, err := b.BuildCompanyList(chatID, user)
	if err != nil {
		return nil
	}
	msg.Text = fmt.Sprintf("%s\n\n%s", text, msg.Text)
	return msg
}
