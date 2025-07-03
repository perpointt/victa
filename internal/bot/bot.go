package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"victa/internal/config"
	"victa/internal/service"
)

// Bot хранит API и ссылку на БД
type Bot struct {
	api             *tgbotapi.BotAPI
	config          config.Config
	UserSvc         *service.UserService
	UserSettingsSvc *service.UserSettingsService
	CompanySvc      *service.CompanyService
}

// NewBot создаёт нового бота
func NewBot(config config.Config, us *service.UserService, uss *service.UserSettingsService, cs *service.CompanyService) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		return nil, err
	}
	return &Bot{api: api, config: config, UserSvc: us, UserSettingsSvc: uss, CompanySvc: cs}, nil
}

// Run запускает цикл обработки команд
func (b *Bot) Run() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.api.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message != nil {
			go b.handleUpdate(update)
		} else if update.CallbackQuery != nil {
			go b.handleCallbackQuery(update.CallbackQuery)
		}
	}
}

// isCallbackWithPrefix возвращает true, если data начинается с указанного префикса.
func isCallbackWithPrefix(data, prefix string) bool {
	return strings.HasPrefix(data, prefix)
}

func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	data := callback.Data

	switch {
	case isCallbackWithPrefix(data, "main_menu"):
		b.HandleMainMenuCallback(callback)

	case isCallbackWithPrefix(data, "company_list"):
		b.HandleListCompaniesCallback(callback)

	case isCallbackWithPrefix(data, "create_company"):
		b.HandleCreateCompanyCallback(callback)

	case isCallbackWithPrefix(data, "cancel_create_company"):
		b.HandleCancelCreateCompanyCallback(callback)

	default:
		AnswerCallback(b.api, callback, "Неизвестное действие.")
	}
}

// AnswerCallback отправляет ответ на callback-запрос.
func AnswerCallback(api *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, text string) {
	answer := tgbotapi.NewCallback(callback.ID, text)
	if _, err := api.Request(answer); err != nil {
		log.Printf("Ошибка ответа на callback: %v", err)
	}
}

// handleUpdate маршрутизует входящие текстовые сообщения и команды.
func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message.IsCommand() {
		b.handleCommand(update.Message)
	} else if update.Message.Text != "" {
		b.handleText(update.Message)
	}
}

func (b *Bot) handleText(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	// Обработка состояния добавления сервера.
	if state, exists := serverAddState[chatID]; exists && state == stateAwaitingAddConfig {
		b.HandleCreateCompany(message)
		return
	}

	b.send(b.newMessage(message.Chat.ID, message.Text))
}

// handleCommand обрабатывает команды (например, /start).
func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		b.HandleStart(message)
	default:
		b.send(b.newMessage(message.Chat.ID, "Неизвестная команда!"))
	}
}

// DeleteMsg удаляет сообщение.
func (b *Bot) deleteMsg(chatID int64, messageID int) {
	delMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	if _, err := b.api.Request(delMsg); err != nil {
		log.Printf("Ошибка удаления сообщения: %v", err)
	}
}

// EditMsg редактирует сообщение с заданным текстом и клавиатурой.
func (b *Bot) editMessage(chatID int64, messageID int, text string, replyMarkup tgbotapi.InlineKeyboardMarkup) {
	editMsg := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, text, replyMarkup)
	editMsg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := b.api.Send(editMsg); err != nil {
		log.Printf("Ошибка редактирования сообщения: %v", err)
	}
}

// newMessage собирает MessageConfig с Markdown-парсингом
func (b *Bot) newMessage(chatID int64, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	return msg
}

// newKeyboardMessage собирает MessageConfig + Inline-клавиатуру
func (b *Bot) newKeyboardMessage(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.MessageConfig {
	msg := b.newMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	return msg
}

// send отправляет любое Chattable (например, MessageConfig)
func (b *Bot) send(c tgbotapi.Chattable) (*tgbotapi.Message, error) {
	msg, err := b.api.Send(c)
	if err != nil {
		log.Printf("send error: %v", err)
		return nil, err
	}

	return &msg, nil
}

// CancelAllOperations отменяет все запущенные процессы, связанные с пользователем.
func (b *Bot) cancelAllOperations(chatID int64) {
	b.CancelCreateCompany(chatID)
}

func (b *Bot) buildCancelButton(data string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData("Отмена", data)
}
