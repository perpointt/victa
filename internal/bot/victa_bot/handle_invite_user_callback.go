package victa_bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleInviteUserCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	tgID := callback.From.ID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверная команда."))
		return
	}

	user, err := b.UserSvc.GetByTgID(ctx, tgID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	if err := b.CompanySvc.CheckAdmin(ctx, user.ID, params.CompanyID); err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	token, err := b.InviteSvc.CreateToken(params.CompanyID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}
	link := fmt.Sprintf("https://t.me/%s?start=%s", b.BotTag, token)
	text := fmt.Sprintf("👤 Ссылка-приглашение (действует 48 ч):\n\n`%s`", link)

	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildCloseButton(),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	b.SendMessage(b.NewKeyboardMessage(chatID, text, keyboard))
}
