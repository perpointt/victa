package victa_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/schema"
	"net/url"
	"strconv"
	"strings"
	"victa/internal/domain"
)

func (b *Bot) AddPendingAppData(chatID int64, data PendingAppData) {
	b.pendingAppData[chatID] = data
}

func (b *Bot) DeletePendingAppData(chatID int64) {
	if _, ok := b.pendingAppData[chatID]; ok {
		delete(b.pendingAppData, chatID)
	}
}

func (b *Bot) AddPendingCompanyID(chatID int64, companyID int64) {
	b.pendingCompanyIDs[chatID] = companyID
}

func (b *Bot) DeletePendingCompanyID(chatID int64) {
	if _, ok := b.pendingCompanyIDs[chatID]; ok {
		delete(b.pendingCompanyIDs, chatID)
	}
}

func (b *Bot) AddPendingAppID(chatID int64, appID int64) {
	b.pendingAppIDs[chatID] = appID
}

func (b *Bot) DeletePendingAppID(chatID int64) {
	if _, ok := b.pendingAppIDs[chatID]; ok {
		delete(b.pendingAppIDs, chatID)
	}
}

func (b *Bot) GetAppDetailMessage(app *domain.App) string {
	const dt = "02.01.2006 15:04:05"

	lines := []string{
		fmt.Sprintf("📱 *%s | %s*", app.Name, app.Slug),
		"",
		fmt.Sprintf("*ID приложения*: `%d`", app.ID),
		fmt.Sprintf("*Создано*: %s", app.CreatedAt.Format(dt)),
		fmt.Sprintf("*Обновлено*: %s", app.UpdatedAt.Format(dt)),
		"",
	}

	// helper: формирует строку «• *Store* — ссылка / нет данных»
	add := func(title string, link *string) {
		if link != nil && *link != "" {
			lines = append(lines, fmt.Sprintf("• *%s* — [%s](%s)", title, "ссылка", *link))
		} else {
			lines = append(lines, fmt.Sprintf("• *%s* — нет данных", title))
		}
	}

	add("App Store", app.AppStoreURL)
	add("Google Play", app.PlayStoreURL)
	add("RuStore", app.RuStoreURL)
	add("Huawei App Gallery", app.AppGalleryURL)

	return strings.Join(lines, "\n")
}

func (b *Bot) GetUserDetailMessage(user *domain.User) string {
	return fmt.Sprintf(
		"*ID пользователя*: `%d`\n*Telegram ID*: `%s`",
		user.ID,
		user.TgID,
	)
}

func (b *Bot) GetCompanyDetailMessage(company *domain.Company) string {
	return fmt.Sprintf(
		"💼 *%s* \n\n*ID компании*: `%d`\n*Создана*: %s\n*Обновлена*: %s",
		company.Name,
		company.ID,
		company.CreatedAt.Format("02.01.2006 15:04:05"),
		company.UpdatedAt.Format("02.01.2006 15:04:05"),
	)
}

func (b *Bot) GetRoleTitle(roleID int64) string {

	switch roleID {
	case 1:
		return "Admin"
	case 2:
		return "Developer"
	default:
		return "Роль #" + strconv.FormatInt(roleID, 10)
	}
}

type CallbackParams struct {
	UserID        int64                `schema:"user_id"`
	CompanyID     int64                `schema:"company_id"`
	AppID         int64                `schema:"app_id"`
	SecretType    domain.SecretType    `schema:"secret_type"`
	AppUpdateType domain.AppUpdateType `schema:"app_update_type"`
}

var schemaDecoder = func() *schema.Decoder {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	return d
}()

func (b *Bot) GetCallbackArgs(data string) (*CallbackParams, error) {
	parts := strings.SplitN(data, "?", 2)
	if len(parts) < 2 {
		return &CallbackParams{}, nil
	}
	vals, err := url.ParseQuery(parts[1])
	if err != nil {
		return nil, err
	}

	var p CallbackParams
	if err := schemaDecoder.Decode(&p, vals); err != nil {
		return nil, err
	}
	return &p, nil
}

func (b *Bot) AddChatState(chatID int64, state ChatState) {
	b.states[chatID] = state
}

func (b *Bot) DeleteChatState(chatID int64) {
	delete(b.states, chatID)
}

func (b *Bot) ClearChatState(chatID int64) {
	b.DeletePendingMessage(chatID)
	b.DeleteChatState(chatID)
	b.DeletePendingCompanyID(chatID)
	b.DeletePendingAppData(chatID)
}

// SendPendingMessage отправляет сообщение и добавляет его ID в очередь для последующего удаления
func (b *Bot) SendPendingMessage(config tgbotapi.MessageConfig) {
	sentMsg := b.SendMessage(config)
	b.pendingMessages[config.ChatID] = append(b.pendingMessages[config.ChatID], sentMsg.MessageID)
}

func (b *Bot) DeletePendingMessage(chatID int64) {
	if msgIDs, ok := b.pendingMessages[chatID]; ok {
		for _, msgID := range msgIDs {
			b.DeleteMessage(chatID, msgID)
		}
		delete(b.pendingMessages, chatID)
	}
}
