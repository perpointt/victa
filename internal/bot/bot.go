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
	InviteSvc  *service.InviteService
	AppSvc     *service.AppService
}

var (
	states            = make(map[int64]ChatState)
	pendingMessages   = make(map[int64][]int)
	pendingCompanyIDs = make(map[int64]int64)
	pendingAppData    = make(map[int64]PendingAppData)
)

type PendingAppData struct {
	ID   int64
	Name string
	Slug string
}

// NewBot создаёт нового бота
func NewBot(
	config config.Config,
	us *service.UserService,
	cs *service.CompanyService,
	is *service.InviteService,
	as *service.AppService,
) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		return nil, err
	}
	return &Bot{
		api:        api,
		config:     config,
		UserSvc:    us,
		CompanySvc: cs,
		InviteSvc:  is,
		AppSvc:     as,
	}, nil
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
		case StateWaitingCreateCompanyName:
			b.HandleCompanyNameCreated(message)
			b.ClearChatState(chatID)
		case StateWaitingUpdateCompanyName:
			b.HandleCompanyNameUpdated(message)
			b.ClearChatState(chatID)
		case StateWaitingUpdateCompanyIntegration:
			b.HandleUpdateCompanyIntegration(message)
			b.ClearChatState(chatID)
		case StateWaitingCreateAppName:
			b.HandleAppNameCreated(message)
		case StateWaitingCreateAppSlug:
			b.HandleAppSlugCreated(message)
			b.ClearChatState(chatID)
		case StateWaitingUpdateAppName:
			b.HandleAppNameUpdated(message)
		case StateWaitingUpdateAppSlug:
			b.HandleAppSlugUpdated(message)
			b.ClearChatState(chatID)
		default:

		}
	} else {
		b.SendMessage(b.NewMessage(message.Chat.ID, message.Text))
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case CommandStart:
		b.HandleStartCommand(message)
	default:
		b.SendMessage(b.NewMessage(message.Chat.ID, "Неизвестная команда!"))
	}
}
