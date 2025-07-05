package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleDeleteAppCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	b.AddChatState(chatID, StateWaitingConfirmDeleteApp)

	msgText := "Подтвердите удаление приложения"
	confirmMessage := b.BuildConfirmMessage(chatID, msgText, fmt.Sprintf("%s?app_id=%v&company_id=%v", CallbackConfirmOperation, params.AppID, params.CompanyID))

	b.SendPendingMessage(confirmMessage)
}

func (b *Bot) HandleConfirmDeleteAppCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	company, err := b.CompanySvc.GetByID(params.CompanyID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при проверке компании."))
		return
	}

	if err := b.AppSvc.Delete(params.AppID); err != nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Не удалось удалить приложение: %v", err)))
		return
	}

	config, err := b.BuildAppList(chatID, company)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Ошибка при построении списка приложений: %v", err)))
		return
	}

	b.SendMessage(*config)
}
