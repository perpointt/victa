package victa_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildMainMenu(chatID int64, user *domain.User) (*tgbotapi.MessageConfig, error) {
	text := fmt.Sprintf("🦊*VICTA*🦊\n\n*Имя пользователя*: %s\n%s", user.Name, b.GetUserDetailMessage(user))

	msg, err := b.BuildCompanyList(chatID, user)
	if err != nil {
		return nil, err
	}
	msg.Text = fmt.Sprintf("%s\n\n%s", text, msg.Text)
	return msg, nil
}
