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

func NewBuildWebhookHandler(jwtSvc *service.JWTService, ci *service.CompanyService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем только POST
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Method not allowed, only POST is accepted", http.StatusMethodNotAllowed)
			return
		}

		// 1) Авторизация
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

		companyIntergration, err := ci.GetCompanyIntegrationByID(companyID)
		if err != nil {
			http.Error(w, "invalid company", http.StatusUnauthorized)
			return
		}

		if companyIntergration == nil {
			http.Error(w, "company not found", http.StatusUnauthorized)
			return
		}

		token := companyIntergration.NotificationBotToken
		if token == nil {
			http.Error(w, "invalid notification token", http.StatusUnauthorized)
			return
		}

		chatID, err := strconv.ParseInt(*companyIntergration.NotificationChatID, 10, 64)
		text := fmt.Sprintf(
			"🚀 Сборка для компании ID `%d` завершена.\n*Build ID:* `%s`",
			companyID, payload.BuildID,
		)

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
