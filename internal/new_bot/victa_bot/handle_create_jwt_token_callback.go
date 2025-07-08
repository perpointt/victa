package victa_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleCreateJwtToken(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при получении данных компании."))
		return
	}

	token, err := b.JwtSvc.GenerateToken(params.CompanyID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при генерации токена."))
		return
	}

	text := fmt.Sprintf("`%s`", token)

	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildCloseButton(),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	b.ClearChatState(chatID)
	b.SendMessage(b.NewKeyboardMessage(chatID, text, keyboard))
}
