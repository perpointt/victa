package notification_bot

import (
	"fmt"
	"strings"
	"time"

	"victa/internal/domain"
)

var ruBuildStatus = map[string]string{
	"publishing": "Сборка завершена",
	"finished":   "Сборка завершена",
	"success":    "Сборка завершена",
	"cancel":     "Сборка отменена",
	"canceled":   "Сборка отменена",
	"failed":     "Ошибка при сборке",
}

var emojiByBuildStatus = map[string]string{
	"publishing": "✅",
	"finished":   "✅",
	"success":    "✅",
	"cancel":     "⚠️",
	"canceled":   "⚠️",
	"failed":     "❌",
}

var buildStepAlias = map[string]string{
	"Set up code signing identities": "Set up code signing",
}

func (bot *Bot) SendDeployNotification(app domain.CodemagicApplication, build domain.CodemagicBuild) {
	text := bot.buildDeployText(app, build)
	_ = bot.SendMessage(bot.NewHtmlMessage(bot.chatID, text))
}

func (bot *Bot) ruBuildStatus(en string) string {
	if v, ok := ruBuildStatus[strings.ToLower(en)]; ok {
		return v
	}
	return en
}

func (bot *Bot) buildStatusEmoji(status string) string {
	if v, ok := emojiByBuildStatus[strings.ToLower(status)]; ok {
		return v
	}
	return "🔹"
}

func (bot *Bot) shortBuildStep(name string) string {
	if v, ok := buildStepAlias[name]; ok {
		return v
	}
	return name
}

func (bot *Bot) buildDeployText(
	app domain.CodemagicApplication,
	build domain.CodemagicBuild,
) string {
	var b strings.Builder
	b.Grow(512)

	fmt.Fprintf(
		&b,
		"<b>🚀 %s | %s %s</b>\n",
		bot.Escape(app.AppName),
		bot.Escape(bot.ruBuildStatus(build.Status)),
		bot.Escape(bot.buildStatusEmoji(build.Status)),
	)

	for _, art := range build.Artefacts {
		if strings.EqualFold(art.Type, "apk") && art.PublicURL != "" {
			fmt.Fprintf(
				&b,
				"\n📦 <b><a href=\"%s\">Скачать APK</a></b>\n",
				bot.Escape(art.PublicURL),
			)
			break
		}
	}

	var duration time.Duration
	if build.FinishedAt.IsZero() {
		duration = time.Since(build.StartedAt)
	} else {
		duration = build.FinishedAt.Sub(build.StartedAt)
	}
	duration = duration.Round(time.Second)
	version := build.Version
	if version == "" {
		version = "Не определена"
	}

	meta := []string{
		fmt.Sprintf("\n<b>• Версия:</b> %s", bot.Escape(version)),
		fmt.Sprintf("<b>• Время сборки:</b> %s", duration.String()),

		fmt.Sprintf("<b>• ID билда:</b> <code>%s</code>", bot.Escape(build.ID)),
		fmt.Sprintf("<b>• Платформы:</b> %s", bot.Escape(strings.Join(build.Config.BuildSettings.Platforms, ", "))),
		fmt.Sprintf("<b>• Версия Flutter:</b> %s", bot.Escape(build.Config.BuildSettings.FlutterVersion)),

		fmt.Sprintf("<b>• Ветка:</b> %s", bot.Escape(build.Commit.Branch)),
		fmt.Sprintf("<b>• Коммит:</b> <code>%s</code>", bot.Escape(build.Commit.CommitMessage)),
		fmt.Sprintf("<b>• Автор коммита:</b> %s", bot.Escape(build.Commit.AuthorName)),
	}

	for _, m := range meta {
		b.WriteString(m + "\n")
	}

	if strings.ToLower(build.Status) != "success" && build.Message != "" {
		fmt.Fprintf(&b, "\n<i>️Текст ошибки:</i>\n<pre>%s</pre>\n", bot.Escape(build.Message))
	}

	if len(build.BuildActions) > 0 {
		b.WriteString("\n<i>Шаги сборки:</i><blockquote expandable>\n")
		for _, act := range build.BuildActions {
			fmt.Fprintf(
				&b,
				"%s %s\n",
				bot.buildStatusEmoji(act.Status),
				bot.Escape(bot.shortBuildStep(act.Name)),
			)
		}
		b.WriteString("</blockquote>")
	}

	buildURL := fmt.Sprintf(
		"https://codemagic.io/app/%s/build/%s",
		bot.Escape(app.ID),
		bot.Escape(build.ID),
	)

	fmt.Fprintf(
		&b,
		"\n\n🔗 <b><a href=\"%s\">Информация о сборке</a></b>\n",
		bot.Escape(buildURL),
	)

	return b.String()
}
