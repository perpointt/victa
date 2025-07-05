package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildMainMenu(chatID int64, user *domain.User) *tgbotapi.MessageConfig {
	text := fmt.Sprintf("🦊*VICTA*🦊\n\n*Имя пользователя*: %s\n%s\n\nВаши компании ⬇️", user.Name, b.GetUserDetailMessage(user))

	msg, err := b.BuildCompanyList(chatID, user)
	if err != nil {
		return nil
	}
	msg.Text = fmt.Sprintf("%s\n\n%s", text, msg.Text)
	return msg
}
