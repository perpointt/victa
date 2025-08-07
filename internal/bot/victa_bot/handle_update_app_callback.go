package victa_bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/url"
	"strings"
	"victa/internal/domain"
)

func (b *Bot) HandleUpdateAppCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверные параметры."))
		return
	}

	switch params.AppUpdateType {
	case domain.AppUpdateName:
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateAppName)
		b.AddPendingAppID(chatID, params.AppID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте новое имя приложения", keyboard))
	case domain.AppUpdateSlug:
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateAppSlug)
		b.AddPendingAppID(chatID, params.AppID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте новый короткий тэг приложения", keyboard))
	case domain.AppUpdateAppStoreURL:
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateAppStoreURL)
		b.AddPendingAppID(chatID, params.AppID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте новую ссылку приложения в App Store", keyboard))
	case domain.AppUpdatePlayStoreURL:
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdatePlayStoreURL)
		b.AddPendingAppID(chatID, params.AppID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте новую ссылку приложения в Play Store", keyboard))
	case domain.AppUpdateRuStoreURL:
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateRuStoreURL)
		b.AddPendingAppID(chatID, params.AppID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте новую ссылку приложения в Ru Store", keyboard))
	case domain.AppUpdateAppGalleryURL:
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateAppGalleryURL)
		b.AddPendingAppID(chatID, params.AppID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте новую ссылку приложения в App Gallery", keyboard))
	}
}

func (b *Bot) HandleUpdateApp(ctx context.Context, message *tgbotapi.Message, updateType domain.AppUpdateType) {
	chatID := message.Chat.ID
	tgID := message.From.ID
	appID := b.pendingAppIDs[chatID]

	app, err := b.AppSvc.GetByID(ctx, appID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	user, err := b.UserSvc.GetByTgID(ctx, tgID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	text := strings.TrimSpace(message.Text)

	switch updateType {
	case domain.AppUpdateName:
		app.Name = text
	case domain.AppUpdateSlug:
		app.Slug = text
	case domain.AppUpdateAppStoreURL,
		domain.AppUpdatePlayStoreURL,
		domain.AppUpdateRuStoreURL,
		domain.AppUpdateAppGalleryURL:

		if err := validateStoreURL(updateType, text); err != nil {
			b.SendErrorMessage(chatID, err)
			return
		}
		switch updateType {
		case domain.AppUpdateAppStoreURL:
			app.AppStoreURL = &text
		case domain.AppUpdatePlayStoreURL:
			app.PlayStoreURL = &text
		case domain.AppUpdateRuStoreURL:
			app.RuStoreURL = &text
		case domain.AppUpdateAppGalleryURL:
			app.AppGalleryURL = &text
		}

	default:
		return
	}

	// сохраняем
	if _, err = b.AppSvc.Update(ctx, app); err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	config := b.BuildAppDetail(ctx, chatID, app, user)

	b.DeleteMessage(chatID, message.MessageID)
	b.ClearChatState(chatID)
	b.SendMessage(config)
}

func validateStoreURL(t domain.AppUpdateType, raw string) error {
	u, err := url.ParseRequestURI(strings.TrimSpace(raw))
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("некорректная ссылка")
	}

	host := strings.ToLower(u.Host)

	switch t {
	case domain.AppUpdateAppStoreURL:
		if !strings.Contains(host, "apps.apple.com") && !strings.Contains(host, "itunes.apple.com") {
			return fmt.Errorf("для App Store ожидается ссылка на apps.apple.com")
		}
	case domain.AppUpdatePlayStoreURL:
		if !strings.Contains(host, "play.google.com") {
			return fmt.Errorf("для Google Play ожидается ссылка на play.google.com")
		}
	case domain.AppUpdateRuStoreURL:
		if !strings.Contains(host, "rustore.ru") {
			return fmt.Errorf("для RuStore ожидается ссылка на rustore.ru")
		}
	case domain.AppUpdateAppGalleryURL:
		if !strings.Contains(host, "appgallery.huawei.com") {
			return fmt.Errorf("для App Gallery ожидается ссылка на appgallery.huawei.com")
		}
	}
	return nil
}
