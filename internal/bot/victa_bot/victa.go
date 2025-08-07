package victa_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"time"
	"victa/internal/bot/bot_common"
	"victa/internal/domain"
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
	pendingAppIDs     map[int64]int64
	pendingAppData    map[int64]PendingAppData
}

type PendingAppData struct {
	ID   int64
	Name string
	Slug string
}

// New создаёт нового бота
func New(
	base *bot_common.BaseBot,
	botTag string,
	us *service.UserService,
	cs *service.CompanyService,
	is *service.InviteService,
	as *service.AppService,
	js *service.JWTService,
) *Bot {
	return &Bot{
		BaseBot:    base,
		BotTag:     botTag,
		UserSvc:    us,
		CompanySvc: cs,
		InviteSvc:  is,
		AppSvc:     as,
		JwtSvc:     js,

		states:            make(map[int64]ChatState),
		pendingMessages:   make(map[int64][]int),
		pendingCompanyIDs: make(map[int64]int64),
		pendingAppIDs:     make(map[int64]int64),
		pendingAppData:    make(map[int64]PendingAppData),
	}
}

const perUpdateTimeout = 10 * time.Second // в конфиг/const

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.BotAPI.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case upd, ok := <-updates:
			if !ok {
				return ctx.Err()
			}

			updCtx, cancel := context.WithTimeout(ctx, perUpdateTimeout)
			b.dispatch(updCtx, upd)
			cancel()
		}
	}
}

func (b *Bot) dispatch(ctx context.Context, upd tgbotapi.Update) {
	switch {
	case upd.Message != nil:
		b.handleMessage(ctx, upd.Message)
	case upd.CallbackQuery != nil:
		b.handleCallback(ctx, upd.CallbackQuery)
	}
}

func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	if msg.IsCommand() {
		b.handleCommand(ctx, msg)
		return
	}

	if msg.Document != nil {
		b.handleDocument(ctx, msg)
		return
	}

	if msg.Text != "" {
		b.handleText(ctx, msg)
		return
	}

	// Остальные типы (фото, стикеры) игнорируем
}

func (b *Bot) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	switch msg.Command() {
	case CommandStart:
		b.HandleStartCommand(ctx, msg)
	default:
		b.SendMessage(b.NewMessage(msg.Chat.ID, "Неизвестная команда!"))
	}
}

func (b *Bot) handleDocument(ctx context.Context, message *tgbotapi.Message) {
	chatID := message.Chat.ID

	if state, exists := b.states[chatID]; exists {
		doc := message.Document
		if doc.FileID == "" {
			b.SendMessage(b.NewMessage(chatID, "Пустой файл"))
			return
		}

		switch state {
		case StateWaitingUploadGoogleJSON:
			b.HandleUpdateCompanySecret(ctx, message, domain.SecretGoogleJSON)
		case StateWaitingUploadAppleP8:
			b.HandleUpdateCompanySecret(ctx, message, domain.SecretAppleP8)
		default:
		}
	}
}

func (b *Bot) handleText(ctx context.Context, message *tgbotapi.Message) {
	chatID := message.Chat.ID

	if state, exists := b.states[chatID]; exists {
		switch state {
		case StateWaitingCreateCompanyName:
			b.HandleCompanyNameCreated(ctx, message)
		case StateWaitingUpdateCompanyName:
			b.HandleCompanyNameUpdated(ctx, message)

		case StateWaitingUpdateCodemagicApiKey:
			b.HandleUpdateCompanySecret(ctx, message, domain.SecretCodemagicApiKey)
		case StateWaitingUpdateAppleKeyID:
			b.HandleUpdateCompanySecret(ctx, message, domain.SecretAppleKeyID)
		case StateWaitingUpdateAppleIssuerID:
			b.HandleUpdateCompanySecret(ctx, message, domain.SecretAppleIssuerID)
		case StateWaitingUpdateNotificationBotToken:
			b.HandleUpdateCompanySecret(ctx, message, domain.SecretNotificationBotToken)
		case StateWaitingUpdateDeployNotificationChatID:
			b.HandleUpdateCompanySecret(ctx, message, domain.SecretDeployNotificationChatID)
		case StateWaitingUpdateIssueNotificationChatID:
			b.HandleUpdateCompanySecret(ctx, message, domain.SecretIssuesNotificationChatID)
		case StateWaitingUpdateErrorNotificationChatID:
			b.HandleUpdateCompanySecret(ctx, message, domain.SecretErrorsNotificationChatID)
		case StateWaitingUpdateVersionNotificationChatID:
			b.HandleUpdateCompanySecret(ctx, message, domain.SecretVersionsNotificationChatID)
		case StateWaitingUpdateReviewsNotificationChatID:
			b.HandleUpdateCompanySecret(ctx, message, domain.SecretReviewsNotificationChatID)

		case StateWaitingUploadGoogleJSON:
			b.DeleteMessage(chatID, message.MessageID)
		case StateWaitingUploadAppleP8:
			b.DeleteMessage(chatID, message.MessageID)

		case StateWaitingCreateAppName:
			b.HandleAppNameCreated(message)
		case StateWaitingCreateAppSlug:
			b.HandleAppSlugCreated(ctx, message)

		case StateWaitingUpdateAppName:
			b.HandleUpdateApp(ctx, message, domain.AppUpdateName)
		case StateWaitingUpdateAppSlug:
			b.HandleUpdateApp(ctx, message, domain.AppUpdateSlug)
		case StateWaitingUpdateAppStoreURL:
			b.HandleUpdateApp(ctx, message, domain.AppUpdateAppStoreURL)
		case StateWaitingUpdatePlayStoreURL:
			b.HandleUpdateApp(ctx, message, domain.AppUpdatePlayStoreURL)
		case StateWaitingUpdateRuStoreURL:
			b.HandleUpdateApp(ctx, message, domain.AppUpdateRuStoreURL)
		case StateWaitingUpdateAppGalleryURL:
			b.HandleUpdateApp(ctx, message, domain.AppUpdateAppGalleryURL)
		default:
		}
	}
}

func (b *Bot) isCallbackWithPrefix(data, prefix string) bool {
	return strings.HasPrefix(data, prefix)
}

func (b *Bot) handleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	b.AnswerCallback(callback, "")

	switch {
	case b.isCallbackWithPrefix(data, CallbackConfirmOperation):
		if state, exists := b.states[chatID]; exists {
			switch state {
			case StateWaitingConfirmDeleteCompany:
				b.HandleConfirmDeleteCompanyCallback(ctx, callback)
				b.ClearChatState(chatID)
			case StateWaitingConfirmDeleteUser:
				b.HandleConfirmDeleteUserCallback(ctx, callback)
				b.ClearChatState(chatID)
			case StateWaitingConfirmDeleteApp:
				b.HandleConfirmDeleteAppCallback(ctx, callback)
				b.ClearChatState(chatID)
			default:
				b.AnswerCallback(callback, "Неизвестное действие.")
			}
		} else {
			b.AnswerCallback(callback, "Неизвестное действие.")
		}
	case b.isCallbackWithPrefix(data, CallbackMainMenu):
		b.ClearChatState(chatID)
		b.HandleMainMenuCallback(ctx, callback)

	case b.isCallbackWithPrefix(data, CallbackClearState):
		b.ClearChatState(chatID)

	case b.isCallbackWithPrefix(data, CallbackDeleteMessage):
		b.HandleDeleteMessageCallback(callback)

	case b.isCallbackWithPrefix(data, CallbackListCompany):
		b.ClearChatState(chatID)
		b.HandleListCompaniesCallback(ctx, callback)
	case b.isCallbackWithPrefix(data, CallbackDetailCompany):
		b.ClearChatState(chatID)
		b.HandleDetailCompanyCallback(ctx, callback)
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
		b.HandleBackToDetailCompanyCallback(ctx, callback)

	case b.isCallbackWithPrefix(data, CallbackCompanyIntegrations):
		b.ClearChatState(chatID)
		b.HandleCompanyIntegrationsCallback(ctx, callback)
	case b.isCallbackWithPrefix(data, CallbackCreateJwtToken):
		b.ClearChatState(chatID)
		b.HandleCreateJwtToken(callback)

	case b.isCallbackWithPrefix(data, CallbackUpdateCompanyIntegrations):
		b.ClearChatState(chatID)
		b.HandleUpdateCompanyIntegrationCallback(callback)

	case b.isCallbackWithPrefix(data, CallbackListUser):
		b.ClearChatState(chatID)
		b.HandleListUsersCallback(ctx, callback)
	case b.isCallbackWithPrefix(data, CallbackInviteUser):
		b.ClearChatState(chatID)
		b.HandleInviteUserCallback(ctx, callback)
	case b.isCallbackWithPrefix(data, CallbackDetailUser):
		b.ClearChatState(chatID)
		b.HandleDetailUserCallback(ctx, callback)
	case b.isCallbackWithPrefix(data, CallbackDeleteUser):
		b.ClearChatState(chatID)
		b.HandleDeleteUserCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackBackToDetailUser):
		b.ClearChatState(chatID)
		b.HandleBackToDetailUserCallback(ctx, callback)

	case b.isCallbackWithPrefix(data, CallbackListApp):
		b.ClearChatState(chatID)
		b.HandleListAppsCallback(ctx, callback)
	case b.isCallbackWithPrefix(data, CallbackCreateApp):
		b.ClearChatState(chatID)
		b.HandleCreateAppCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackUpdateApp):
		b.ClearChatState(chatID)
		b.HandleUpdateAppCallback(callback)
	case b.isCallbackWithPrefix(data, CallbackDetailApp):
		b.ClearChatState(chatID)
		b.HandleDetailAppCallback(ctx, callback)
	case b.isCallbackWithPrefix(data, CallbackDeleteApp):
		b.ClearChatState(chatID)
		b.HandleDeleteAppCallback(callback)

	default:
		b.AnswerCallback(callback, "Неизвестное действие.")
	}

}
