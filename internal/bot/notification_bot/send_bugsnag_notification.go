package notification_bot

import (
	"fmt"
	"strings"
	"time"
	"victa/internal/domain"
)

func (bot *Bot) SendBugsnagNotification(w domain.BugsnagWebhook) {
	text := bot.buildErrorText(w)
	_ = bot.SendMessage(bot.NewHtmlMessage(bot.chatID, text))
}

func (bot *Bot) buildErrorText(w domain.BugsnagWebhook) string {
	var b strings.Builder
	b.Grow(512)

	projectName := bot.Escape(w.Project.Name)
	projectURL := bot.Escape(w.Project.URL)
	device := w.Error.Device
	app := w.Error.App

	b.WriteString(fmt.Sprintf("⚠️ <b><a href=\"%s\">%s</a> | %s</b>\n",
		projectURL, projectName, bot.buildErrorTitle(w)))

	meta := []string{
		fmt.Sprintf("\n<b>• Версия приложения:</b> %s+%s", bot.Escape(app.Version), bot.Escape(app.VersionCode)),
		fmt.Sprintf("<b>• Платформа:</b> %s", bot.Escape(app.Type)),
		fmt.Sprintf("<b>• Статус:</b> %s", bot.Escape(bot.buildErrorStatus(w))),
		fmt.Sprintf("<b>• Происшествия:</b> %d", w.Error.Occurrences),
		fmt.Sprintf("<b>• Обработана:</b> %v", !w.Error.Unhandled),
		fmt.Sprintf("<b>• User ID:</b> <code>%v</code>", w.Error.UserID),
	}

	for _, m := range meta {
		b.WriteString(m + "\n")
	}

	b.WriteString(fmt.Sprintf("\n<i>️Текст ошибки:</i>\n<pre>%s</pre>",
		bot.Escape(w.Error.Message)))

	b.WriteString("\n\n<i>Информация об устройстве:</i>\n<blockquote expandable>")

	deviceMeta := []string{
		fmt.Sprintf("<b>• Устройство:</b> %s %s", bot.Escape(device.Manufacturer), bot.Escape(device.Model)),
		fmt.Sprintf("<b>• OS:</b> %s %s", bot.Escape(device.OSName), bot.Escape(device.OSVersion)),
		fmt.Sprintf("<b>• Locale:</b> %s", bot.Escape(device.Locale)),
		fmt.Sprintf("<b>• Orientation:</b> %s", bot.Escape(device.Orientation)),
		fmt.Sprintf("<b>• Battery:</b> %.0f %% (Charging: %v)", device.BatteryLevel*100, device.Charging),
		fmt.Sprintf("<b>• RAM:</b> %s свободно из %s",
			bot.formatBytes(device.FreeMemory), bot.formatBytes(device.TotalMemory)),
		fmt.Sprintf("<b>• Disk:</b> %s свободно", bot.formatBytes(device.FreeDisk)),
		fmt.Sprintf("<b>• Jailbreak:</b> %v", device.JailBroken),
		fmt.Sprintf("<b>• Время на устройстве:</b> %s", device.Time.Format(time.RFC3339)),
	}

	for _, m := range deviceMeta {
		b.WriteString(m + "\n")
	}

	b.WriteString("</blockquote>")

	b.WriteString("\n\n<i>StackTrace:</i>\n<blockquote expandable>")

	for _, ex := range w.Error.Exceptions {
		if len(ex.StackTrace) == 0 {
			continue
		}

		for _, frame := range ex.StackTrace {
			b.WriteString(bot.Escape(
				fmt.Sprintf("%s:%s — %s\n\n",
					frame.File,
					frame.LineNumber,
					frame.Method,
				)))
		}
	}

	b.WriteString("</blockquote>")

	b.WriteString(fmt.Sprintf("\n\n🔗 <b><a href=\"%s\">Информация об ошибке</a></b>",
		bot.Escape(w.Error.URL)))

	return b.String()
}

func (bot *Bot) formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func (bot *Bot) buildErrorTitle(w domain.BugsnagWebhook) string {
	switch w.Trigger.Type {
	case "firstException":
		return "Новая ошибка"
	case "errorEventFrequency":
		return "Ошибка возникает часто"
	case "reopened":
		return "Повторное открытие ошибки"
	case "projectSpiking":
		return "Всплеск исключений в проекте"
	case "errorStateManualChange":
		return "Статус ошибки изменено вручную"
	default:
		return w.Trigger.Type
	}
}

func (bot *Bot) buildErrorStatus(w domain.BugsnagWebhook) string {
	switch w.Error.Status {
	case "open":
		return "Открыта"
	case "fixed":
		return "Исправлена"
	case "snoozed":
		return "Отложена"
	case "ignored":
		return "Игнорирована"
	default:
		return w.Error.Status
	}
}
