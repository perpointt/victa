package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// HandleDetailCompanyCallback обрабатывает нажатие кнопки детальной информации о компании
// и редактирует текущее сообщение, выводя информацию по выбранной компании.
func (b *Bot) HandleDetailCompanyCallback(cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID

	// Подтверждаем получение callback, чтобы убрать «часики»
	if _, err := b.api.Request(tgbotapi.NewCallback(cb.ID, "")); err != nil {
		log.Printf("callback answer error: %v", err)
	}

	// Парсим ID компании из callback.Data: "<action>:<id>"
	idPtr, err := b.GetIdFromCallback(cb.Data)
	if err != nil || idPtr == nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}
	companyID := *idPtr

	// Получаем компанию из БД
	company, err := b.CompanySvc.GetById(companyID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при получении данных компании."))
		return
	}
	if company == nil {
		b.SendMessage(b.NewMessage(chatID, "Компания не найдена."))
		return
	}

	// Формируем текст с информацией о компании
	text := fmt.Sprintf(
		"*%s* (ID: %d)\nСоздана: %s\nОбновлена: %s",
		company.Name,
		company.ID,
		company.CreatedAt.Format("02 Jan 2006 15:04"),
		company.UpdatedAt.Format("02 Jan 2006 15:04"),
	)

	rows := tgbotapi.NewInlineKeyboardRow(
		b.BuildCloseButton(CallbackDeleteMessage),
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows)

	b.ClearChatState(chatID)
	b.SendMessage(b.NewKeyboardMessage(chatID, text, keyboard))
}
