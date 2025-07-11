package notification_bot

import (
	"fmt"
	"strings"
	"victa/internal/domain"
)

var ruIssueStatus = map[string]string{
	"opened": "Открыта",
	"closed": "Закрыта",
}

func (bot *Bot) SendIssueNotification(issue domain.GitlabWebhook) {
	text := bot.buildIssueText(issue)
	_ = bot.SendMessage(bot.NewHtmlMessage(bot.chatID, text))
}

func (bot *Bot) buildIssueText(issue domain.GitlabWebhook) string {
	var b strings.Builder
	b.Grow(512)

	obj := bot.getIssueObject(issue)
	projectName := bot.Escape(issue.Project.Name)
	projectURL := bot.Escape(issue.Project.Homepage)
	issueURL := bot.Escape(obj.URL)
	issueIID := fmt.Sprintf("%d", obj.IID)

	fmt.Fprintf(&b,
		"<b><a href=\"%s\">%s</a> | <a href=\"%s\">#%s</a></b>\n\n",
		projectURL, projectName, issueURL, issueIID,
	)

	fmt.Fprintf(&b,
		"• Задача: <b>%s</b>\n",
		bot.Escape(obj.Title),
	)
	fmt.Fprintf(&b,
		"• Статус: <b>%s</b>\n\n",
		bot.Escape(bot.ruIssueStatus(obj.State)),
	)

	switch issue.ObjectKind {
	case "issue":
		switch issue.ObjectAttributes.Action {
		case "open", "reopen":
			b.WriteString("🚀 <b>Задача открыта</b>")
		case "close":
			b.WriteString("✅ <b>Задача закрыта</b>")
		case "update":
			b.WriteString("🔄 <b>Задача обновлена</b>")
		default:
			fmt.Fprintf(&b, "<b>%s</b>", issue.ObjectAttributes.Action)
		}

		fmt.Fprintf(&b,
			"<i> by %s</i>\n\n",
			bot.Escape(issue.User.Name),
		)

	case "note":
		switch issue.ObjectAttributes.Action {
		case "create":
			b.WriteString("💬 <b>Новый комментарий</b>")
		case "update":
			b.WriteString("💬 <b>Комментарий отредактирован</b>")
		default:
			fmt.Fprintf(&b, "<b>%s</b>", issue.ObjectAttributes.Action)
		}

		fmt.Fprintf(&b,
			"<i> by %s</i>\n",
			bot.Escape(issue.User.Name),
		)

		commentURL := issue.ObjectAttributes.URL

		fmt.Fprintf(&b,
			"🔗 <a href=\"%s\">Ссылка на комментарий</a>\n\n",
			commentURL,
		)

		comment := issue.ObjectAttributes.Description
		if comment != "" {
			fmt.Fprintf(&b,
				"<i>️Текст комментария:</i>\n<blockquote expandable>%s</blockquote>\n\n",
				bot.MarkdownToHTML(comment),
			)
		}
	}

	desc := bot.MarkdownToHTML(obj.Description)
	if desc != "" {
		fmt.Fprintf(&b,
			"<i>Описание задачи:</i>\n<blockquote expandable>%s</blockquote>\n",
			desc,
		)
	}

	return b.String()
}

func (bot *Bot) getIssueObject(issue domain.GitlabWebhook) domain.Attributes {
	if issue.ObjectKind == "note" {
		return issue.Issue
	}
	return issue.ObjectAttributes
}

func (bot *Bot) ruIssueStatus(en string) string {
	if v, ok := ruIssueStatus[strings.ToLower(en)]; ok {
		return v
	}
	return en
}
