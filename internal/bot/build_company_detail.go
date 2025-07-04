package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildCompanyDetail(chatID int64, company *domain.Company) tgbotapi.MessageConfig {
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
		tgbotapi.NewInlineKeyboardButtonData("Сотрудники", fmt.Sprintf("%s:%v", CallbackListUser, company.ID)),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Интеграции", CallbackCompanyIntegrations),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildDeleteButton(fmt.Sprintf("%s:%v", CallbackDeleteCompany, company.ID)),
		b.BuildEditButton(fmt.Sprintf("%s:%v", CallbackUpdateCompany, company.ID)),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildCloseButton(CallbackDeleteMessage),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	message := b.NewKeyboardMessage(chatID, text, keyboard)
	return message
}
