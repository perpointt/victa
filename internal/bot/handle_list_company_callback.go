package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleListCompaniesCallback обрабатывает нажатие кнопки «Компании»
// и редактирует текущее сообщение, выводя каждую компанию отдельной inline-кнопкой
func (b *Bot) HandleListCompaniesCallback(cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID
	messageID := cb.Message.MessageID
	tgID := cb.From.ID

	// Закрываем «часики» в кнопке
	if _, err := b.api.Request(tgbotapi.NewCallback(cb.ID, "")); err != nil {
		log.Printf("answer callback error: %v", err)
	}

	// Проверяем пользователя
	user, err := b.UserSvc.FindByTgID(tgID)
	if err != nil {
		b.send(b.newMessage(chatID, "Ошибка при поиске пользователя."))
		return
	}
	if user == nil {
		b.send(b.newMessage(chatID, "Сначала зарегистрируйтесь через /start."))
		return
	}

	msgCfg := b.BuildCompanyList(chatID, user)
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
