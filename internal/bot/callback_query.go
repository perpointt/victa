package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (b *Bot) isCallbackWithPrefix(data, prefix string) bool {
	return strings.HasPrefix(data, prefix)
}

func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	data := callback.Data

	switch {
	case b.isCallbackWithPrefix(data, CallbackMainMenu):
		b.HandleMainMenuCallback(callback)

	case b.isCallbackWithPrefix(data, CallbackClearState):
		b.HandleClearStateCallback(callback)

	case b.isCallbackWithPrefix(data, CallbackDeleteMessage):
		b.HandleDeleteMessageCallback(callback)

	case b.isCallbackWithPrefix(data, CallbackListCompany):
		b.HandleListCompaniesCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackDetailCompany):
		b.HandleDetailCompanyCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackCreateCompany):
		b.HandleCreateCompanyCallback(callback)
	default:
		b.AnswerCallback(callback, "Неизвестное действие.")
	}
}
