package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *Bot) BuildCloseButton(data string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData("❌ Закрыть", data)
}

func (b *Bot) BuildBackButton(data string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", data)
}

func (b *Bot) BuildCancelButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", CallbackClearState)
}

func (b *Bot) BuildConfirmButton(data string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить", data)
}

func (b *Bot) BuildDeleteButton(data string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить", data)
}

func (b *Bot) BuildEditButton(data string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData("✏️ Изменить", data)
}
