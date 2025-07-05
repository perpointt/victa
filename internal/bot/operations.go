package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/schema"
	"log"
	"net/url"
	"strconv"
	"strings"
	"victa/internal/domain"
)

func (b *Bot) AddPendingAppData(chatID int64, data PendingAppData) {
	pendingAppData[chatID] = data
}

func (b *Bot) DeletePendingAppData(chatID int64) {
	if _, ok := pendingAppData[chatID]; ok {
		delete(pendingAppData, chatID)
	}
}

func (b *Bot) AddPendingCompanyID(chatID int64, companyID int64) {
	pendingCompanyIDs[chatID] = companyID
}

func (b *Bot) DeletePendingCompanyID(chatID int64) {
	if _, ok := pendingCompanyIDs[chatID]; ok {
		delete(pendingCompanyIDs, chatID)
	}
}

func (b *Bot) GetAppDetailMessage(app *domain.App) string {
	return fmt.Sprintf(
		"📱 *%s | %s* \n\n*ID приложения*: `%d`\n*Создано*: %s\n*Обновлено*: %s",
		app.Name,
		app.Slug,
		app.ID,
		app.CreatedAt.Format("02.01.2006 15:04:05"),
		app.UpdatedAt.Format("02.01.2006 15:04:05"),
	)
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
	UserID    int64 `schema:"user_id"`
	CompanyID int64 `schema:"company_id"`
	AppID     int64 `schema:"app_id"`
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
	states[chatID] = state
}

func (b *Bot) DeleteChatState(chatID int64) {
	delete(states, chatID)
}

func (b *Bot) ClearChatState(chatID int64) {
	b.DeletePendingMessage(chatID)
	b.DeleteChatState(chatID)
	b.DeletePendingCompanyID(chatID)
	b.DeletePendingAppData(chatID)
}

// SendPendingMessage отправляет сообщение и добавляет его ID в очередь для последующего удаления
func (b *Bot) SendPendingMessage(config tgbotapi.MessageConfig) {
	sentMsg, err := b.SendMessage(config)
	if err != nil {
		return
	}

	pendingMessages[config.ChatID] = append(pendingMessages[config.ChatID], sentMsg.MessageID)
}

func (b *Bot) DeletePendingMessage(chatID int64) {
	if msgIDs, ok := pendingMessages[chatID]; ok {
		for _, msgID := range msgIDs {
			_ = b.DeleteMessage(chatID, msgID)
		}
		delete(pendingMessages, chatID)
	}
}

func (b *Bot) AnswerCallback(callback *tgbotapi.CallbackQuery, text string) error {
	answer := tgbotapi.NewCallback(callback.ID, text)
	if _, err := b.api.Request(answer); err != nil {
		log.Printf("Ошибка ответа на callback: %v", err)
		return err
	}

	return nil
}

func (b *Bot) SendMessage(config tgbotapi.MessageConfig) (*tgbotapi.Message, error) {
	msg, err := b.api.Send(config)

	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
		return nil, err
	}

	return &msg, nil
}

func (b *Bot) EditMessage(messageID int, config tgbotapi.MessageConfig) (*tgbotapi.Message, error) {
	text := config.Text
	replyMarkup, ok := config.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)

	if !ok {
		return nil, fmt.Errorf("expected InlineKeyboardMarkup, got %T", config.ReplyMarkup)
	}

	editMsg := tgbotapi.NewEditMessageTextAndMarkup(config.ChatID, messageID, text, replyMarkup)
	editMsg.ParseMode = config.ParseMode

	msg, err := b.api.Send(editMsg)

	if err != nil {
		log.Printf("Ошибка редактирования сообщения: %v", err)
		return nil, err
	}

	return &msg, nil
}

func (b *Bot) DeleteMessage(chatID int64, messageID int) error {
	delMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	if _, err := b.api.Request(delMsg); err != nil {
		log.Printf("Ошибка удаления сообщения: %v", err)
		return err
	}

	return nil
}

func (b *Bot) NewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

func (b *Bot) NewKeyboardMessage(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.MessageConfig {
	msg := b.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboard
	return msg
}
