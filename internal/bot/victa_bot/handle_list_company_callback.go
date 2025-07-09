package victa_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleListCompaniesCallback(cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID
	tgID := cb.From.ID

	user, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		return
	}
	if user == nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Сначала зарегистрируйтесь через /%v.", CommandStart)))
		return
	}

	message, err := b.BuildCompanyList(chatID, user)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Ошибка при построении списка компаний: %v", err)))
		return
	}

	b.EditMessage(messageID, *message)
}
