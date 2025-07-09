package victa_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"victa/internal/bot/bot_common"
	"victa/internal/service"
)

// Bot хранит API и ссылку на БД
type Bot struct {
	*bot_common.BaseBot
	BotTag     string
	UserSvc    *service.UserService
	CompanySvc *service.CompanyService
	InviteSvc  *service.InviteService
	AppSvc     *service.AppService
	JwtSvc     *service.JWTService

	states            map[int64]ChatState
	pendingMessages   map[int64][]int
	pendingCompanyIDs map[int64]int64
	pendingAppData    map[int64]PendingAppData
}

type PendingAppData struct {
	ID   int64
	Name string
	Slug string
}

// NewBot создаёт нового бота
func NewBot(
	base *bot_common.BaseBot,
	us *service.UserService,
	cs *service.CompanyService,
	is *service.InviteService,
	as *service.AppService,
	js *service.JWTService,
) *Bot {
	return &Bot{
		BaseBot:    base,
		UserSvc:    us,
		CompanySvc: cs,
		InviteSvc:  is,
		AppSvc:     as,
		JwtSvc:     js,

		states:            make(map[int64]ChatState),
		pendingMessages:   make(map[int64][]int),
		pendingCompanyIDs: make(map[int64]int64),
		pendingAppData:    make(map[int64]PendingAppData),
	}
}

func (b *Bot) Run() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.BotAPI.GetUpdatesChan(updateConfig)
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

	if state, exists := b.states[chatID]; exists {
		switch state {
		case StateWaitingCreateCompanyName:
			b.HandleCompanyNameCreated(message)
		case StateWaitingUpdateCompanyName:
			b.HandleCompanyNameUpdated(message)
		case StateWaitingUpdateCompanyIntegration:
			b.HandleUpdateCompanyIntegration(message)
		case StateWaitingCreateAppName:
			b.HandleAppNameCreated(message)
		case StateWaitingCreateAppSlug:
			b.HandleAppSlugCreated(message)
		case StateWaitingUpdateAppName:
			b.HandleAppNameUpdated(message)
		case StateWaitingUpdateAppSlug:
			b.HandleAppSlugUpdated(message)
		default:
		}
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

func (b *Bot) isCallbackWithPrefix(data, prefix string) bool {
	return strings.HasPrefix(data, prefix)
}

func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	b.AnswerCallback(callback, "")

	switch {
	case b.isCallbackWithPrefix(data, CallbackConfirmOperation):
		if state, exists := b.states[chatID]; exists {
			switch state {
			case StateWaitingConfirmDeleteCompany:
				b.HandleConfirmDeleteCompanyCallback(callback)
				b.ClearChatState(chatID)
			case StateWaitingConfirmDeleteUser:
				b.HandleConfirmDeleteUserCallback(callback)
				b.ClearChatState(chatID)
			case StateWaitingConfirmDeleteApp:
				b.HandleConfirmDeleteAppCallback(callback)
				b.ClearChatState(chatID)
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
	case b.isCallbackWithPrefix(data, CallbackCreateJwtToken):
		b.ClearChatState(chatID)
		b.HandleCreateJwtToken(callback)
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

	case b.isCallbackWithPrefix(data, CallbackListApp):
		b.ClearChatState(chatID)
		b.HandleListAppsCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackCreateApp):
		b.ClearChatState(chatID)
		b.HandleCreateAppCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackUpdateApp):
		b.ClearChatState(chatID)
		b.HandleUpdateAppCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackDetailApp):
		b.ClearChatState(chatID)
		b.HandleDetailAppCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackDeleteApp):
		b.ClearChatState(chatID)
		b.HandleDeleteAppCallback(callback)

	default:
		b.AnswerCallback(callback, "Неизвестное действие.")
	}

}
