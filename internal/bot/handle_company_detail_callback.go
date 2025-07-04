package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleDetailCompanyCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID

	idPtr, err := b.GetIdFromCallback(callback.Data)
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
		"*%s* (ID: %d)\n\nСоздана: %s\nОбновлена: %s",
		company.Name,
		company.ID,
		company.CreatedAt.Format("02 Jan 2006 15:04"),
		company.UpdatedAt.Format("02 Jan 2006 15:04"),
	)

	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Приложения", CallbackListApp),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Сотрудники", CallbackListUser),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Интеграции", CallbackCompanyIntegrations),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildDeleteButton(fmt.Sprintf("%s:%v", CallbackDeleteCompany, companyID)),
		b.BuildEditButton(fmt.Sprintf("%s:%v", CallbackUpdateCompany, companyID)),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildCloseButton(CallbackDeleteMessage),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	b.ClearChatState(chatID)
	b.SendMessage(b.NewKeyboardMessage(chatID, text, keyboard))
}
