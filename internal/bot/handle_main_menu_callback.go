package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleMainMenuCallback(cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID
	// … проверяем пользователя …

	// Строим основное меню
	user, _ := b.UserSvc.FindByTgID(cb.From.ID)
	msgCfg := b.BuildMainMenu(chatID, user)
	if msgCfg == nil {
		return
	}

	// Извлекаем текст и клавиатуру
	text := msgCfg.Text
	kb, ok := msgCfg.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)
	if !ok {
		return
	}

	b.editMessage(chatID, messageID, text, kb)
}
