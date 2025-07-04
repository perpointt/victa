package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleDeleteCompanyCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	idPtr, err := b.GetIdFromCallback(callback.Data)
	if err != nil || idPtr == nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}
	companyID := *idPtr

	b.AddChatState(chatID, StateWaitingConfirmDeleteCompany)

	msgText := "Подтвердите удаление компании"
	confirmMessage := b.BuildConfirmMessage(chatID, msgText, fmt.Sprintf("%s:%v", CallbackConfirmOperation, companyID))

	b.SendPendingMessage(confirmMessage)
}

func (b *Bot) HandleConfirmDeleteCompanyCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	tgID := callback.From.ID
	idPtr, err := b.GetIdFromCallback(callback.Data)
	if err != nil || idPtr == nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}
	companyID := *idPtr

	user, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при проверке пользователя."))
		return
	}
	if user == nil {
		b.SendMessage(b.NewMessage(chatID, "Сначала зарегистрируйтесь через /start."))
		return
	}

	if err := b.CompanySvc.Delete(companyID, user.ID); err != nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Не удалось удалить компанию: %v", err)))
		return
	}

	mainMenuMsg := b.BuildMainMenu(chatID, user)
	if mainMenuMsg == nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Ошибка при построении главного меню: %v", err)))
		return
	}

	b.SendMessage(*mainMenuMsg)
}
