package notification_bot

import (
	"bytes"
	"github.com/yuin/goldmark"
	"html"
	"regexp"
	"strconv"
	"strings"
	"victa/internal/bot/bot_common"
)

// Bot хранит API и ссылку на БД
type Bot struct {
	*bot_common.BaseBot
	chatID int64
}

// NewBot создаёт нового бота
func NewBot(
	base *bot_common.BaseBot,
	chatIDString string,
) (*Bot, error) {
	chatID, err := strconv.ParseInt(chatIDString, 10, 64)
	if err != nil {
		return nil, err
	}

	return &Bot{
		BaseBot: base,
		chatID:  chatID,
	}, nil
}

func (bot *Bot) Escape(s string) string { return html.EscapeString(s) }

// MarkdownToHTML конвертит Markdown в HTML.
func (bot *Bot) MarkdownToHTML(md string) string {
	// 1) Markdown → HTML
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return md
	}
	htmlStr := buf.String()

	// 2) Распаковываем html-сущности (&quot; → ")
	htmlStr = html.UnescapeString(htmlStr)

	// 3) Заменяем <strong>/<em> → <b>/<i>, и <p>…</p> на пустое + двойной \n
	htmlStr = strings.ReplaceAll(htmlStr, "<strong>", "<b>")
	htmlStr = strings.ReplaceAll(htmlStr, "</strong>", "</b>")
	htmlStr = strings.ReplaceAll(htmlStr, "<em>", "<i>")
	htmlStr = strings.ReplaceAll(htmlStr, "</em>", "</i>")
	htmlStr = strings.ReplaceAll(htmlStr, "<p>", "")
	htmlStr = strings.ReplaceAll(htmlStr, "</p>", "\n\n")

	// 4) Убираем все теги <…>, но:
	//    – если это <b>, </b>, <i>, </i>, <u>, </u>, <code>, </code>, <pre>, </pre>, or <a href="…">/</a>,
	//      оставляем;
	//    – если это <br> или <br/>, заменяем на \n;
	//    – всё остальное выкидываем.
	re := regexp.MustCompile(`<[^>]+>`)
	result := re.ReplaceAllStringFunc(htmlStr, func(tag string) string {
		t := strings.ToLower(tag)
		switch {
		case t == "<b>" || t == "</b>",
			t == "<i>" || t == "</i>",
			t == "<u>" || t == "</u>",
			t == "<code>" || t == "</code>",
			t == "<pre>" || t == "</pre>":
			return tag
		case strings.HasPrefix(t, `<a `) && strings.HasSuffix(t, `>`):
			return tag
		case t == "</a>":
			return "</a>"
		case t == "<br>" || t == "<br/>":
			return "\n"
		default:
			return ""
		}
	})

	// 5) Убираем лишние пробелы у концов и возвращаем
	return strings.TrimSpace(result)
}
