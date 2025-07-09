package victa_bot

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"reflect"
	"strings"
	"victa/internal/domain"
)

// HandleCompanyIntegrationsCallback обрабатывает нажатие кнопки «Интеграции»
func (b *Bot) HandleCompanyIntegrationsCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	params, err := b.GetCallbackArgs(callback.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверные параметры."))
		return
	}
	companyID := params.CompanyID

	company, err := b.CompanySvc.GetByID(companyID)
	if err != nil {
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		return
	}

	config, err := b.BuildCompanyIntegrationsDetail(chatID, company)
	if err != nil {
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		return
	}
	b.EditMessage(messageID, *config)
}

// BuildIntegrationTemplate автоматически собирает JSON-шаблон
// для всех полей CompanyIntegration (кроме company_id).
func (b *Bot) BuildIntegrationTemplate(ci *domain.CompanyIntegration) (string, error) {
	// 1) Разворачиваем nil-поинтерн в zero-value struct
	inst := domain.CompanyIntegration{}
	if ci != nil {
		inst = *ci
	}

	// 2) Готовим отражения типа и значения
	rt := reflect.TypeOf(inst)
	rv := reflect.ValueOf(inst)

	// 3) Собираем map[tag]value
	m := make(map[string]string, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if tag == "" || tag == "company_id" {
			continue
		}
		fv := rv.Field(i)
		var val string
		// все поля кроме company_id — это *string
		if fv.IsNil() {
			val = ""
		} else {
			val = fv.Elem().String()
		}
		m[tag] = val
	}

	// 4) Красиво маршалим
	bytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
