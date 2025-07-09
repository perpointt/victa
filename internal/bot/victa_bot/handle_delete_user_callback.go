package victa_bot

import (
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

func (b *Bot) HandleConfirmDeleteUserCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	tgID := callback.From.ID
	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	err = b.UserSvc.DeleteFromCompany(params.UserID, params.CompanyID)
	if err != nil {
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		return
	}

	company, err := b.CompanySvc.GetByID(params.CompanyID)
	if err != nil {
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		return
	}

	userList, err := b.BuildUserList(chatID, tgID, company)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Ошибка при построении списка пользователей: %v", err)))
		return
	}

	b.SendMessage(*userList)
}
