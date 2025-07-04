package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/config"
	"victa/internal/service"
)

// Bot хранит API и ссылку на БД
type Bot struct {
	api        *tgbotapi.BotAPI
	config     config.Config
	UserSvc    *service.UserService
	CompanySvc *service.CompanyService
}

var (
	states                   = make(map[int64]ChatState)
	pendingMessages          = make(map[int64]int)
	pendingUpdateCompanyData = make(map[int64]int64)
)

// NewBot создаёт нового бота
func NewBot(config config.Config, us *service.UserService, cs *service.CompanyService) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		return nil, err
	}
	return &Bot{api: api, config: config, UserSvc: us, CompanySvc: cs}, nil
}

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

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message.IsCommand() {
		b.handleCommand(update.Message)
	} else if update.Message.Text != "" {
		b.handleText(update.Message)
	}
}

func (b *Bot) handleText(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	if state, exists := states[chatID]; exists {
		switch state {
		case StateWaitingCreateCompany:
			b.HandleCreateCompany(message)
			b.ClearChatState(chatID)
			return
		case StateWaitingUpdateCompany:
			b.HandleUpdateCompany(message)
			b.ClearChatState(chatID)
			return
		default:

		}
	} else {
		b.SendMessage(b.NewMessage(message.Chat.ID, message.Text))
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case CommandStart:
		b.HandleStart(message)
	default:
		b.SendMessage(b.NewMessage(message.Chat.ID, "Неизвестная команда!"))
	}
}
