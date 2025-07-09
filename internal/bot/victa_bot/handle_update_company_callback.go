package victa_bot

import (
	"fmt"
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

func (b *Bot) HandleCompanyNameUpdated(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	tgID := message.From.ID

	user, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		return
	}

	companyID := b.pendingCompanyIDs[chatID]

	_, err = b.CompanySvc.Update(companyID, message.Text, user.ID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Не удалось обновить компанию: %v", err)))
		return
	}

	menu, err := b.BuildMainMenu(chatID, user)
	if err != nil {
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		return
	}

	b.SendMessage(*menu)
}
