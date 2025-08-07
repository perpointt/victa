package victa_bot

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"strings"
	"victa/internal/domain"
)

func (b *Bot) HandleUpdateCompanyIntegrationCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверные параметры."))
		return
	}

	// helper: убираем «notification» из SecretType → короче callback_data
	shortType := func(st domain.SecretType) string {
		return strings.ReplaceAll(string(st), "notification", "")
	}

	switch string(params.SecretType) {
	case string(domain.SecretCodemagicApiKey):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateCodemagicApiKey)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Codemagic API Key", keyboard))
	case shortType(domain.SecretNotificationBotToken):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateNotificationBotToken)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Notification Bot Token", keyboard))
	case shortType(domain.SecretDeployNotificationChatID):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateDeployNotificationChatID)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Deploy Notification Chat ID", keyboard))
	case shortType(domain.SecretIssuesNotificationChatID):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateIssueNotificationChatID)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Issue Notification Chat ID", keyboard))
	case shortType(domain.SecretErrorsNotificationChatID):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateErrorNotificationChatID)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Error Notification Chat ID", keyboard))
	case shortType(domain.SecretVersionsNotificationChatID):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateVersionNotificationChatID)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Versions Notification Chat ID", keyboard))
	case shortType(domain.SecretReviewsNotificationChatID):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateReviewsNotificationChatID)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Reviews Notification Chat ID", keyboard))
	case shortType(domain.SecretAppleP8):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUploadAppleP8)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте AppleP8", keyboard))
	case shortType(domain.SecretAppleKeyID):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateAppleKeyID)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Apple Key ID", keyboard))
	case shortType(domain.SecretAppleIssuerID):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)

		b.AddChatState(chatID, StateWaitingUpdateAppleIssuerID)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Apple Issuer ID", keyboard))
	case shortType(domain.SecretGoogleJSON):
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				b.BuildCancelButton(),
			),
		)
		b.AddChatState(chatID, StateWaitingUploadGoogleJSON)
		b.AddPendingCompanyID(chatID, params.CompanyID)
		b.SendPendingMessage(b.NewKeyboardMessage(chatID, "Отправьте Google JSON", keyboard))
	}
}

func (b *Bot) HandleUpdateCompanySecret(
	ctx context.Context,
	message *tgbotapi.Message,
	secretType domain.SecretType,
) {
	chatID := message.Chat.ID
	companyID := b.pendingCompanyIDs[chatID]

	company, err := b.CompanySvc.GetByID(ctx, companyID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	short := func(st domain.SecretType) string {
		return strings.ReplaceAll(string(st), "notification", "")
	}

	switch string(secretType) {

	case short(domain.SecretAppleP8), short(domain.SecretGoogleJSON):
		if message.Document == nil {
			b.SendErrorMessage(chatID, fmt.Errorf("пришлите файл документом"))
			return
		}

		data, err := b.downloadFile(message.Document.FileID)
		if err != nil {
			b.SendErrorMessage(chatID, fmt.Errorf("скачивание файла: %v", err))
			return
		}

		if secretType == domain.SecretGoogleJSON && !json.Valid(data) {
			b.SendErrorMessage(chatID, fmt.Errorf("файл не похож на валидный JSON"))
			return
		}
		if secretType == domain.SecretAppleP8 && !strings.HasPrefix(string(data), "-----BEGIN PRIVATE KEY-----") {
			b.SendErrorMessage(chatID, fmt.Errorf("ожидается .p8-файл (BEGIN PRIVATE KEY)"))
			return
		}

		if _, err := b.CompanySvc.CreateBinarySecret(ctx, company.ID, secretType, data); err != nil {
			b.SendErrorMessage(chatID, err)
			return
		}

	default:
		if _, err :=
			b.CompanySvc.CreateTextSecret(ctx, company.ID, secretType, message.Text); err != nil {
			b.SendErrorMessage(chatID, err)
			return
		}
	}

	config, err := b.BuildCompanyIntegrationsDetail(ctx, chatID, company)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.DeleteMessage(chatID, message.MessageID)
	b.ClearChatState(chatID)
	b.SendMessage(*config)
}

func (b *Bot) downloadFile(fileID string) ([]byte, error) {
	file, err := b.BotAPI.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return nil, fmt.Errorf("getFile: %w", err)
	}

	url := file.Link(b.BotAPI.Token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	const maxSize = 5 << 20
	limited := io.LimitReader(resp.Body, maxSize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if len(data) > maxSize {
		return nil, fmt.Errorf("файл больше %d MB", maxSize>>20)
	}
	return data, nil
}
