package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleInviteUserCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	tgID := callback.From.ID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	user, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при проверке пользователя."))
		return
	}

	if err := b.CompanySvc.CheckAdmin(user.ID, params.CompanyID); err != nil {
		b.SendMessage(b.NewMessage(chatID, fmt.Sprintf("Не удалось создать приглашение: %v", err)))
		return
	}

	token := b.InviteSvc.CreateToken(params.CompanyID)
	link := fmt.Sprintf("https://t.me/%s?start=%s", b.config.TelegramBotName, token)
	text := fmt.Sprintf("Ссылка-приглашение (действует 48 ч):\n```%s```", link)

	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildCloseButton(CallbackDeleteMessage),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	b.SendMessage(b.NewKeyboardMessage(chatID, text, keyboard))
}
