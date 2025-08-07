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
		has := map[domain.AppUpdateType]bool{
			domain.AppUpdateName:          app.Name != "",
			domain.AppUpdateSlug:          app.Slug != "",
			domain.AppUpdateAppStoreURL:   app.AppStoreURL != nil && *app.AppStoreURL != "",
			domain.AppUpdatePlayStoreURL:  app.PlayStoreURL != nil && *app.PlayStoreURL != "",
			domain.AppUpdateRuStoreURL:    app.RuStoreURL != nil && *app.RuStoreURL != "",
			domain.AppUpdateAppGalleryURL: app.AppGalleryURL != nil && *app.AppGalleryURL != "",
		}

		mark := func(ok bool) string {
			if ok {
				return "🟢"
			}
			return "🔴"
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s Slug", mark(has[domain.AppUpdateSlug])),
				fmt.Sprintf("%s?app_id=%d&app_update_type=%s",
					CallbackUpdateApp,
					app.ID,
					domain.AppUpdateSlug),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s Name", mark(has[domain.AppUpdateName])),
				fmt.Sprintf("%s?app_id=%d&app_update_type=%s",
					CallbackUpdateApp,
					app.ID,
					domain.AppUpdateName),
			),
		))

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s App Store URL", mark(has[domain.AppUpdateAppStoreURL])),
				fmt.Sprintf("%s?app_id=%d&app_update_type=%s",
					CallbackUpdateApp,
					app.ID,
					domain.AppUpdateAppStoreURL),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s Play Store URL", mark(has[domain.AppUpdatePlayStoreURL])),
				fmt.Sprintf("%s?app_id=%d&app_update_type=%s",
					CallbackUpdateApp,
					app.ID,
					domain.AppUpdatePlayStoreURL),
			),
		))

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s Ru Store URL", mark(has[domain.AppUpdateRuStoreURL])),
				fmt.Sprintf("%s?app_id=%d&app_update_type=%s",
					CallbackUpdateApp,
					app.ID,
					domain.AppUpdateRuStoreURL),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s App Gallery URL", mark(has[domain.AppUpdateAppGalleryURL])),
				fmt.Sprintf("%s?app_id=%d&app_update_type=%s",
					CallbackUpdateApp,
					app.ID,
					domain.AppUpdateAppGalleryURL),
			),
		))

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			b.BuildDeleteButton(fmt.Sprintf("%v?app_id=%d&company_id=%d", CallbackDeleteApp, app.ID, app.CompanyID)),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		b.BuildCloseButton(),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	message := b.NewKeyboardMessage(chatID, text, keyboard)
	return message
}
