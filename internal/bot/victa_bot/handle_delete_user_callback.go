package victa_bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleDeleteUserCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	b.AddChatState(chatID, StateWaitingConfirmDeleteUser)

	msgText := "Подтвердите удаление пользователя из компании"
	confirmMessage := b.BuildConfirmMessage(chatID, msgText, fmt.Sprintf("%s?company_id=%v&user_id=%v", CallbackConfirmOperation, params.CompanyID, params.UserID))

	b.SendPendingMessage(confirmMessage)
}

func (b *Bot) HandleConfirmDeleteUserCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	tgID := callback.From.ID
	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	err = b.UserSvc.DeleteFromCompany(ctx, params.UserID, params.CompanyID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	company, err := b.CompanySvc.GetByID(ctx, params.CompanyID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	userList, err := b.BuildUserList(ctx, chatID, tgID, company)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.SendMessage(*userList)
}
