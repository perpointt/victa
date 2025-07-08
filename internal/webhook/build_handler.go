package webhook

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"strconv"
	"strings"

	"victa/internal/service"
)

// BuildWebhookPayload описывает тело запроса на ваш webhook.
type BuildWebhookPayload struct {
	BuildID string `json:"build_id"`
}

func NewBuildWebhookHandler(jwtSvc *service.JWTService, ci *service.CompanyService, cm *service.CodemagicService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем только POST
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Method not allowed, only POST is accepted", http.StatusMethodNotAllowed)
			return
		}

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		companyID, err := jwtSvc.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		var payload BuildWebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		integration, err := ci.GetCompanyIntegrationByID(companyID)
		if err != nil {
			http.Error(w, "invalid company", http.StatusUnauthorized)
			return
		}

		if integration == nil {
			http.Error(w, "company not found", http.StatusUnauthorized)
			return
		}

		buildResp, err := cm.GetBuildByID(payload.BuildID, *integration.CodemagicAPIKey)
		if err != nil {
			http.Error(w, "failed to fetch build info", http.StatusInternalServerError)
			return
		}

		token := integration.NotificationBotToken
		if token == nil {
			http.Error(w, "invalid notification token", http.StatusUnauthorized)
			return
		}

		chatID, err := strconv.ParseInt(*integration.NotificationChatID, 10, 64)

		app := buildResp.Application
		build := buildResp.Build

		// 1) Собираем заголовок и базовые поля
		lines := []string{
			fmt.Sprintf("🚀 *%s* (%s)", app.AppName, app.ID),
			"",
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
			fmt.Sprintf("• Ветка: `%s`", build.Commit.Branch),
			"",
		}

		// 2) Добавляем секцию buildActions
		if len(build.BuildActions) > 0 {
			lines = append(lines, "⚙️ *Шаги сборки:*")
			for _, act := range build.BuildActions {
				// можно эмодзи по статусу
				emoji := map[string]string{"success": "✅", "failed": "❌"}[act.Status]
				if emoji == "" {
					emoji = "🔸"
				}
				lines = append(lines,
					fmt.Sprintf("%s %s — `%s`", emoji, act.Name, act.Status),
				)
			}
			lines = append(lines, "")
		}

		// 3) Сообщение от Codemagic
		if build.Message != "" {
			lines = append(lines,
				fmt.Sprintf("💬 *Сообщение:* %s", build.Message),
				"",
			)
		}

		// 4) Артефакты
		if len(build.Artefacts) > 0 {
			lines = append(lines, "📦 *Артефакты:*")
			for _, art := range build.Artefacts {
				lines = append(lines,
					fmt.Sprintf("• [%s](%s) — `%s`", art.Name, art.URL, art.Type),
				)
			}
		}

		// 5) Собираем весь текст
		text := strings.Join(lines, "\n")

		api, err := tgbotapi.NewBotAPI(*token)

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = tgbotapi.ModeMarkdown

		_, err = api.Send(msg)

		if err != nil {
			http.Error(w, "failed to send notification", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
