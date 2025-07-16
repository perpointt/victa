package victa_bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) BuildAppDetail(ctx context.Context, chatID int64, app *domain.App, user *domain.User) tgbotapi.MessageConfig {
	text := b.GetAppDetailMessage(app)

	var rows [][]tgbotapi.InlineKeyboardButton

	err := b.CompanySvc.CheckAdmin(ctx, user.ID, app.CompanyID)

	if err == nil {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			b.BuildDeleteButton(fmt.Sprintf("%v?app_id=%d&company_id=%d", CallbackDeleteApp, app.ID, app.CompanyID)),
			b.BuildEditButton(fmt.Sprintf("%v?app_id=%d&company_id=%d", CallbackUpdateApp, app.ID, app.CompanyID)),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildCloseButton(),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	message := b.NewKeyboardMessage(chatID, text, keyboard)
	return message
}
