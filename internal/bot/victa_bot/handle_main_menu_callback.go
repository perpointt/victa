package victa_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleMainMenuCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	user, err := b.UserSvc.GetByTgID(ctx, callback.From.ID)

	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	if user == nil {
		b.SendMessage(b.NewMessage(chatID, "Сначала зарегистрируйтесь через /start."))
		return
	}

	menu, err := b.BuildMainMenu(ctx, chatID, user)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.ClearChatState(chatID)
	b.EditMessage(messageID, *menu)
}
