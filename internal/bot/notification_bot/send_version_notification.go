package notification_bot

import (
	"fmt"
	"strings"
	"victa/internal/domain"
)

// SendVersionNotification отправляет уведомление о новой версии.
func (bot *Bot) SendVersionNotification(info domain.ReleaseInfo) {
	text := bot.buildNewVersionText(info)
	_ = bot.SendMessage(bot.NewHtmlMessage(bot.chatID, text))
}

// buildNewVersionText формирует текст уведомления о новой версии.
func (bot *Bot) buildNewVersionText(info domain.ReleaseInfo) string {
	var b strings.Builder
	b.Grow(256)

	var title string
	switch info.Store {
	case domain.StoreGooglePlay:
		title = "🤖 Google Play"
	case domain.StoreAppStore:
		title = "🍏 App Store"
	default:
		title = string(info.Store)
	}

	b.WriteString(fmt.Sprintf("🔔 <b>Новая версия в %s</b>\n", bot.Escape(title)))

	versionLine := fmt.Sprintf("• <b>Версия:</b> %s+%d", bot.Escape(info.Semantic), info.Code)
	b.WriteString(versionLine)
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("• <b>App ID:</b> <code>%s</code>\n", bot.Escape(info.AppID)))

	if info.BundleID != "" {
		b.WriteString(fmt.Sprintf("• <b>Bundle ID:</b> <code>%s</code>\n", bot.Escape(info.BundleID)))
	}

	if !info.ReleasedAt.IsZero() {
		b.WriteString(fmt.Sprintf("• <b>Дата релиза:</b> %s\n",
			bot.Escape(info.ReleasedAt.Format("02.01.2006 15:04"))))
	}

	return b.String()
}
