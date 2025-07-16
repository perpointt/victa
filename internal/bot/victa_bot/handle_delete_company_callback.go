package victa_bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleDeleteCompanyCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	b.AddChatState(chatID, StateWaitingConfirmDeleteCompany)

	msgText := "Подтвердите удаление компании"
	confirmMessage := b.BuildConfirmMessage(chatID, msgText, fmt.Sprintf("%s?company_id=%v", CallbackConfirmOperation, params.CompanyID))

	b.SendPendingMessage(confirmMessage)
}

func (b *Bot) HandleConfirmDeleteCompanyCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	tgID := callback.From.ID
	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	user, err := b.UserSvc.GetByTgID(ctx, tgID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}
	if user == nil {
		b.SendMessage(b.NewMessage(chatID, "Сначала зарегистрируйтесь через /start."))
		return
	}

	if err := b.CompanySvc.Delete(ctx, params.CompanyID, user.ID); err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	menu, err := b.BuildMainMenu(ctx, chatID, user)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.SendMessage(*menu)
}
