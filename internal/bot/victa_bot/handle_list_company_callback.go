package victa_bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleListCompaniesCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID
	tgID := callback.From.ID

	user, err := b.UserSvc.GetByTgID(ctx, tgID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}
	if user == nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Сначала зарегистрируйтесь через /%v.", CommandStart)))
		return
	}

	message, err := b.BuildCompanyList(ctx, chatID, user)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Ошибка при построении списка компаний: %v", err)))
		return
	}

	b.EditMessage(messageID, *message)
}
