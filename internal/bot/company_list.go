package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildCompanyList(chatID int64, user *domain.User) *tgbotapi.MessageConfig {
	// Загружаем компании пользователя
	companies, err := b.CompanySvc.GetAllByUserId(user.ID)
	if err != nil {
		b.send(b.newMessage(chatID, "Ошибка при получении списка компаний."))
		return nil
	}

	// Текст сообщения
	text := "Ваши компании:"
	if len(companies) == 0 {
		text = "У вас ещё нет компаний."
	}

	// Формируем клавиатуру: по одной кнопке на компанию
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range companies {
		// callbackData можно настроить как нужно, например "company_<id>"
		cbData := fmt.Sprintf("select_company_%d", c.ID)
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(c.Name, cbData),
			),
		)
	}

	// Добавляем кнопку создания компании
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ Создать компанию", "create_company"),
	))
	// И кнопку «Назад»
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "main_menu"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := b.newKeyboardMessage(chatID, text, keyboard)
	return &msg
}
