package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleListCompaniesCallback(cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID
	tgID := cb.From.ID

	user, err := b.UserSvc.FindByTgID(tgID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при поиске пользователя."))
		return
	}
	if user == nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Сначала зарегистрируйтесь через /%v.", CommandStart)))
		return
	}

	config := b.BuildCompanyList(chatID, user)
	if config == nil {
		return
	}

	b.EditMessage(messageID, *config)
}
