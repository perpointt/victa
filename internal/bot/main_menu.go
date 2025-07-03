package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

// BuildMainMenu отправляет главное меню, показывая «Пользователи» и «Компании» только админу
func (b *Bot) BuildMainMenu(chatID int64, user *domain.User) *tgbotapi.MessageConfig {
	isAdmin := user.TgId == b.config.TelegramAdminUserId
	// получаем настроенную активную компанию
	settings, err := b.UserSettingsSvc.FindByUserID(user.ID)
	if err != nil {
		return nil
	}

	var activeCompany string
	if settings.ActiveCompanyId != nil {
		activeCompany = fmt.Sprintf("%d", *settings.ActiveCompanyId)
	} else {
		activeCompany = "Нет компании"
	}

	text := fmt.Sprintf(
		"Привет, *%s*!\n\n*Активная компания:* _%s_",
		user.Name, activeCompany,
	)

	// Сборка строк клавиатуры
	var rows [][]tgbotapi.InlineKeyboardButton

	// Всегда доступно
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📱 Приложения", "menu_apps"),
	))

	// Только для администратора
	if isAdmin {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("👥 Пользователи", "menu_users"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🏢 Компании", "company_list"),
			),
		)
	}

	// Всегда доступно
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "menu_settings"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// Отправляем через утилиты
	msg := b.newKeyboardMessage(chatID, text, keyboard)
	return &msg
}
