package victa_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleUpdateCompanyCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	b.AddChatState(chatID, StateWaitingUpdateCompanyName)
	b.AddPendingCompanyID(chatID, params.CompanyID)

	msgText := "Отправьте название компании"
	cancelButton := b.BuildCancelButton()
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(cancelButton))

	b.SendPendingMessage(b.NewKeyboardMessage(chatID, msgText, keyboard))
}

func (b *Bot) HandleCompanyNameUpdated(ctx context.Context, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	tgID := message.From.ID

	user, err := b.UserSvc.GetByTgID(ctx, tgID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	companyID := b.pendingCompanyIDs[chatID]

	_, err = b.CompanySvc.Update(ctx, companyID, message.Text, user.ID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	menu, err := b.BuildMainMenu(ctx, chatID, user)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.ClearChatState(chatID)

	b.SendMessage(*menu)
}
