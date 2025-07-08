package notification_bot

import (
	"fmt"
	"strings"
	"victa/internal/domain"
)

func (bot *Bot) SendNewNotification(app domain.CodemagicApplication, build domain.CodemagicBuild) error {
	text := bot.buildCodemagicText(app, build)
	_, err := bot.SendMessage(bot.NewMessage(bot.chatID, text))
	if err != nil {
		return err
	}

	return nil
}

// buildCodemagicText собирает Markdown-отчёт по ответу Codemagic.
func (bot *Bot) buildCodemagicText(app domain.CodemagicApplication, build domain.CodemagicBuild) string {
	lines := []string{
		fmt.Sprintf("🚀 *%s* (%s)\n", app.AppName, app.ID),
		fmt.Sprintf("• *Статус:* `%s`", build.Status),
		fmt.Sprintf("• *Build ID:* `%s`", build.ID),
		fmt.Sprintf("• *Начата:* %s", build.StartedAt.Format("02.01.2006 15:04:05")),
		fmt.Sprintf("• *Завершена:* %s", build.FinishedAt.Format("02.01.2006 15:04:05")),
		"",
		fmt.Sprintf("💻 *Workflow:* %s", build.Config.Name),
		fmt.Sprintf("• Flutter %s  |  Платформы: %s",
			build.Config.BuildSettings.FlutterVersion,
			strings.Join(build.Config.BuildSettings.Platforms, ", "),
		),
		"",
		fmt.Sprintf("🔀 *Коммит:* `%s`", build.Commit.CommitMessage),
		fmt.Sprintf("• _%s_", build.Commit.AuthorName),
		fmt.Sprintf("• Ветка: `%s`\n", build.Commit.Branch),
	}

	if len(build.BuildActions) > 0 {
		lines = append(lines, "⚙️ *Шаги сборки:*")
		for _, act := range build.BuildActions {
			emoji := map[string]string{"success": "✅", "failed": "❌"}[act.Status]
			if emoji == "" {
				emoji = "🔸"
			}
			lines = append(lines, fmt.Sprintf("%s %s — `%s`", emoji, act.Name, act.Status))
		}
		lines = append(lines, "")
	}

	if build.Message != "" {
		lines = append(lines, fmt.Sprintf("💬 *Сообщение:* %s\n", build.Message))
	}

	if len(build.Artefacts) > 0 {
		lines = append(lines, "📦 *Артефакты:*")
		for _, art := range build.Artefacts {
			lines = append(lines, fmt.Sprintf("• [%s](%s) — `%s`", art.Name, art.URL, art.Type))
		}
	}

	return strings.Join(lines, "\n")
}
