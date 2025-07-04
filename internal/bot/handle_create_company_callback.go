package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (b *Bot) HandleCreateCompanyCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	b.AddChatState(chatID, StateWaitingCreateCompany)

	msgText := "Отправьте название компании"
	cancelButton := b.BuildCancelButton()
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(cancelButton))

	b.SendPendingMessage(b.NewKeyboardMessage(chatID, msgText, keyboard))
}

func (b *Bot) HandleCreateCompany(message *tgbotapi.Message) {
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

	_, err = b.CompanySvc.Create(message.Text, user.ID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при создании компании."))
		log.Fatalf(err.Error())
		return
	}

	config := b.BuildMainMenu(chatID, user)
	if config == nil {
		return
	}

	b.SendMessage(*config)
}
