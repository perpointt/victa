package victa_bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleCreateAppCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	params, err := b.GetCallbackArgs(callback.Data)

	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	msgText := "Отправьте название приложения"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		b.BuildCancelButton(),
	))

	b.AddPendingCompanyID(chatID, params.CompanyID)
	b.AddChatState(chatID, StateWaitingCreateAppName)

	b.SendPendingMessage(b.NewKeyboardMessage(chatID, msgText, keyboard))
}

func (b *Bot) HandleAppNameCreated(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	msgText := "Отправьте короткий тэг приложения"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		b.BuildCancelButton(),
	))

	data := b.pendingAppData[chatID]
	data.Name = message.Text

	b.AddPendingAppData(chatID, data)
	b.AddChatState(chatID, StateWaitingCreateAppSlug)

	b.SendPendingMessage(b.NewKeyboardMessage(chatID, msgText, keyboard))
}

func (b *Bot) HandleAppSlugCreated(ctx context.Context, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	companyID := b.pendingCompanyIDs[chatID]
	tgID := message.From.ID

	data := b.pendingAppData[chatID]
	data.Slug = message.Text

	user, err := b.UserSvc.GetByTgID(ctx, tgID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	app, err := b.AppSvc.Create(ctx, companyID, data.Name, data.Slug)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	config := b.BuildAppDetail(ctx, chatID, app, user)

	b.ClearChatState(chatID)

	b.SendMessage(config)
}
