package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) BuildConfirmMessage(chatID int64, title string, data string) tgbotapi.MessageConfig {
	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildConfirmButton(data),
		b.BuildCancelButton(),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := b.NewKeyboardMessage(chatID, title, keyboard)
	return msg
}
