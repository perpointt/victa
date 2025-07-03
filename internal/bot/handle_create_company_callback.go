package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const stateAwaitingAddConfig = "awaiting_add_config"

// Глобальные переменные для работы с состоянием добавления сервера.
var (
	serverAddState      = make(map[int64]string) // chatID -> состояние процесса
	pendingAddServerMsg = make(map[int64]int)    // chatID -> ID сообщения с инструкцией
)

func (b *Bot) HandleCreateCompanyCallback(cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID

	// Отмена других процессов.
	b.cancelAllOperations(chatID)

	serverAddState[chatID] = stateAwaitingAddConfig

	msgText := "Отправьте название компании"
	cancelButton := b.buildCancelButton("cancel_create_company")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(cancelButton))

	sentMsg, err := b.send(b.newKeyboardMessage(chatID, msgText, keyboard))
	if err == nil {
		pendingAddServerMsg[chatID] = sentMsg.MessageID
	}

	AnswerCallback(b.api, cb, "")
}

// HandleCreateCompany обрабатывает сообщение с конфигурацией сервера.
func (b *Bot) HandleCreateCompany(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	tgID := message.From.ID

	// Проверяем пользователя
	user, err := b.UserSvc.FindByTgID(tgID)
	if err != nil {
		b.send(b.newMessage(chatID, "Ошибка при поиске пользователя."))
		return
	}
	if user == nil {
		b.send(b.newMessage(chatID, "Сначала зарегистрируйтесь через /start."))
		return
	}

	_, err = b.CompanySvc.Create(message.Text, user.ID)
	if err != nil {
		b.send(b.newMessage(chatID, "Ошибка при создании компании."))
		log.Fatalf(err.Error())
		return
	}
	b.cancelAllOperations(chatID)

	b.send(b.BuildCompanyList(chatID, user))
}

// CancelCreateCompany отменяет процесс добавления сервера и очищает связанные с ним состояния.
func (b *Bot) CancelCreateCompany(chatID int64) {
	if msgID, ok := pendingAddServerMsg[chatID]; ok {
		b.deleteMsg(chatID, msgID)
		delete(pendingAddServerMsg, chatID)
	}

	// Независимо от состояния, удаляем запись из serverAddState.
	delete(serverAddState, chatID)
}

// HandleCancelCreateCompanyCallback обрабатывает callback для отмены добавления сервера.
func (b *Bot) HandleCancelCreateCompanyCallback(callback *tgbotapi.CallbackQuery) {
	b.CancelCreateCompany(callback.Message.Chat.ID)
	AnswerCallback(b.api, callback, "Добавление сервера отменено.")
}
