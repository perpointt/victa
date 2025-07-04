package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleUpdateCompanyCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	idPtr, err := b.GetIdFromCallback(callback.Data)
	if err != nil || idPtr == nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}
	companyID := *idPtr

	b.AddChatState(chatID, StateWaitingUpdateCompany)
	b.AddPendingUpdateCompanyData(chatID, companyID)

	msgText := "Отправьте название компании"
	cancelButton := b.BuildCancelButton()
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(cancelButton))

	b.SendPendingMessage(b.NewKeyboardMessage(chatID, msgText, keyboard))
}

func (b *Bot) HandleUpdateCompany(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	tgID := message.From.ID

	user, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при поиске пользователя."))
		return
	}
	if user == nil {
		b.SendMessage(b.NewMessage(chatID, "Сначала зарегистрируйтесь через /start."))
		return
	}

	companyID := pendingUpdateCompanyData[chatID]

	_, err = b.CompanySvc.Update(companyID, message.Text, user.ID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Не удалось обновить компанию: %v", err)))
		return
	}

	config := b.BuildMainMenu(chatID, user)
	if config == nil {
		return
	}

	b.SendMessage(*config)
}

func (b *Bot) AddPendingUpdateCompanyData(chatID int64, companyId int64) {
	pendingUpdateCompanyData[chatID] = companyId
}

func (b *Bot) DeletePendingUpdateCompanyData(chatID int64) {
	if _, ok := pendingUpdateCompanyData[chatID]; ok {
		delete(pendingUpdateCompanyData, chatID)
	}
}
