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
	case b.isCallbackWithPrefix(data, "main_menu"):
		b.HandleMainMenuCallback(callback)
	case b.isCallbackWithPrefix(data, "company_list"):
		b.HandleListCompaniesCallback(callback)
	case b.isCallbackWithPrefix(data, "create_company"):
		b.HandleCreateCompanyCallback(callback)
	case b.isCallbackWithPrefix(data, "cancel_create_company"):
		b.HandleCancelCreateCompanyCallback(callback)
	default:
		b.AnswerCallback(callback, "Неизвестное действие.")
	}
}
