package notification_bot

import (
	"fmt"
	"html"
	"strings"
	"time"

	"victa/internal/domain"
)

var ruStatus = map[string]string{
	"failed":   "Ошибка при сборке",
	"cancel":   "Сборка отменена",
	"finished": "Сборка завершена",
}

var emojiByStat = map[string]string{
	"success":  "✅",
	"finished": "✅",
	"canceled": "⚠️",
	"failed":   "❌",
}

var stepAlias = map[string]string{
	"Set up code signing identities": "Set up code signing",
}

func (bot *Bot) SendNewNotification(app domain.CodemagicApplication, build domain.CodemagicBuild) error {
	text := bot.buildCodemagicText(app, build)
	_, err := bot.SendMessage(bot.NewHtmlMessage(bot.chatID, text))
	if err != nil {
		return err
	}

	return nil
}

func esc(s string) string { return html.EscapeString(s) }

func ruBuildStatus(en string) string {
	if v, ok := ruStatus[strings.ToLower(en)]; ok {
		return v
	}
	return en
}

func emoji(status string) string {
	if v, ok := emojiByStat[strings.ToLower(status)]; ok {
		return v
	}
	return "🔹"
}

func shortStep(name string) string {
	if v, ok := stepAlias[name]; ok {
		return v
	}
	return name
}

func (bot *Bot) buildCodemagicText(
	app domain.CodemagicApplication,
	build domain.CodemagicBuild,
) string {
	var b strings.Builder
	b.Grow(512)

	fmt.Fprintf(
		&b,
		"<b>🚀 %s | %s %s</b>\n",
		esc(app.AppName),
		esc(ruBuildStatus(build.Status)),
		esc(emoji(build.Status)),
	)

	for _, art := range build.Artefacts {
		if strings.EqualFold(art.Type, "apk") && art.PublicURL != "" {
			fmt.Fprintf(
				&b,
				"\n📦 <b><a href=\"%s\">Скачать APK</a></b>\n",
				esc(art.PublicURL),
			)
			break
		}
	}

	duration := build.FinishedAt.Sub(build.StartedAt).Round(time.Second)
	version := build.Version
	if version == "" {
		version = "Не определена"
	}

	meta := []string{
		fmt.Sprintf("\n<b>Версия:</b> %s", esc(version)),
		fmt.Sprintf("<b>Время сборки:</b> %s", esc(duration.String())),

		fmt.Sprintf("\n<b>ID билда:</b> <code>%s</code>", esc(build.ID)),
		fmt.Sprintf("<b>Платформы:</b> %s", esc(strings.Join(build.Config.BuildSettings.Platforms, ", "))),
		fmt.Sprintf("<b>Версия Flutter:</b> %s", esc(build.Config.BuildSettings.FlutterVersion)),

		fmt.Sprintf("\n<b>Ветка:</b> %s", esc(build.Commit.Branch)),
		fmt.Sprintf("<b>Коммит:</b> %s", esc(build.Commit.CommitMessage)),
		fmt.Sprintf("<b>Автор коммита:</b> %s", esc(build.Commit.AuthorName)),
	}

	for _, m := range meta {
		b.WriteString(m + "\n")
	}

	if strings.ToLower(build.Status) != "success" && build.Message != "" {
		fmt.Fprintf(&b, "\n<pre>%s</pre>\n", esc(build.Message))
	}

	if len(build.BuildActions) > 0 {
		b.WriteString("\n<blockquote expandable>⚙️ <b>Шаги сборки</b>:\n")
		for _, act := range build.BuildActions {
			fmt.Fprintf(
				&b,
				"%s %s\n",
				emoji(act.Status),
				esc(shortStep(act.Name)),
			)
		}
		b.WriteString("</blockquote>")
	}

	buildURL := fmt.Sprintf(
		"https://codemagic.io/app/%s/build/%s",
		esc(app.ID),
		esc(build.ID),
	)

	fmt.Fprintf(
		&b,
		"\n\n🔗 <b><a href=\"%s\">Информация о сборке</a></b>\n",
		esc(buildURL),
	)

	return b.String()
}
