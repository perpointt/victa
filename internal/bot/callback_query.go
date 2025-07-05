package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (b *Bot) isCallbackWithPrefix(data, prefix string) bool {
	return strings.HasPrefix(data, prefix)
}

func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	b.AnswerCallback(callback, "")

	switch {
	case b.isCallbackWithPrefix(data, CallbackConfirmOperation):
		if state, exists := states[chatID]; exists {
			switch state {
			case StateWaitingConfirmDeleteCompany:
				b.HandleConfirmDeleteCompanyCallback(callback)
				b.ClearChatState(chatID)
				return
			case StateWaitingConfirmDeleteUser:
				b.HandleConfirmDeleteUserCallback(callback)
				b.ClearChatState(chatID)
				return
			default:
				b.AnswerCallback(callback, "Неизвестное действие.")
			}
		} else {
			b.AnswerCallback(callback, "Неизвестное действие.")
		}
	case b.isCallbackWithPrefix(data, CallbackMainMenu):
		b.ClearChatState(chatID)
		b.HandleMainMenuCallback(callback)

	case b.isCallbackWithPrefix(data, CallbackClearState):
		b.ClearChatState(chatID)

	case b.isCallbackWithPrefix(data, CallbackDeleteMessage):
		b.HandleDeleteMessageCallback(callback)

	case b.isCallbackWithPrefix(data, CallbackListCompany):
		b.ClearChatState(chatID)
		b.HandleListCompaniesCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackDetailCompany):
		b.ClearChatState(chatID)
		b.HandleDetailCompanyCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackCreateCompany):
		b.ClearChatState(chatID)
		b.HandleCreateCompanyCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackUpdateCompany):
		b.ClearChatState(chatID)
		b.HandleUpdateCompanyCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackDeleteCompany):
		b.ClearChatState(chatID)
		b.HandleDeleteCompanyCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackBackToDetailCompany):
		b.ClearChatState(chatID)
		b.HandleBackToDetailCompanyCallback(callback)

	case b.isCallbackWithPrefix(data, CallbackCompanyIntegrations):
		b.ClearChatState(chatID)
		b.HandleCompanyIntegrationsCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackUpdateCompanyIntegrations):
		b.ClearChatState(chatID)
		b.HandleUpdateCompanyIntegrationCallback(callback)

	case b.isCallbackWithPrefix(data, CallbackListUser):
		b.ClearChatState(chatID)
		b.HandleListUsersCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackInviteUser):
		b.ClearChatState(chatID)
		b.HandleInviteUserCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackDetailUser):
		b.ClearChatState(chatID)
		b.HandleDetailUserCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackDeleteUser):
		b.ClearChatState(chatID)
		b.HandleDeleteUserCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackBackToDetailUser):
		b.ClearChatState(chatID)
		b.HandleBackToDetailUserCallback(callback)

	default:
		b.AnswerCallback(callback, "Неизвестное действие.")
	}

}
