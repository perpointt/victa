package victa_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (b *Bot) HandleUpdateAppCallback(callback *tgbotapi.CallbackQuery) {
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

	b.AddPendingAppData(chatID, PendingAppData{ID: params.AppID})

	b.AddPendingCompanyID(chatID, params.CompanyID)
	b.AddChatState(chatID, StateWaitingUpdateAppName)

	b.SendPendingMessage(b.NewKeyboardMessage(chatID, msgText, keyboard))
}

func (b *Bot) HandleAppNameUpdated(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	msgText := "Отправьте короткий тэг приложения"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		b.BuildCancelButton(),
	))

	data := b.pendingAppData[chatID]
	data.Name = message.Text

	b.AddPendingAppData(chatID, data)
	b.AddChatState(chatID, StateWaitingUpdateAppSlug)

	b.SendPendingMessage(b.NewKeyboardMessage(chatID, msgText, keyboard))
}

func (b *Bot) HandleAppSlugUpdated(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	tgID := message.From.ID

	data := b.pendingAppData[chatID]
	data.Slug = message.Text

	user, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	app, err := b.AppSvc.Update(data.ID, data.Name, data.Slug)
	if err != nil {
		log.Printf(fmt.Sprintf("%v", err.Error()))

		b.SendErrorMessage(chatID, err)
		return
	}

	config := b.BuildAppDetail(chatID, app, user)

	b.SendMessage(config)
}
