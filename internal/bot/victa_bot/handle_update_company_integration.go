package victa_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleUpdateCompanyIntegrationCallback(cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID

	params, err := b.GetCallbackArgs(cb.Data)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Неверные параметры."))
		return
	}
	companyID := params.CompanyID

	ci, err := b.CompanySvc.GetCompanyIntegrationByID(companyID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	tmpl, err := b.BuildIntegrationTemplate(ci)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	msgText := fmt.Sprintf(
		"Отправьте обновленный JSON с данными для интеграции:\n\n```json\n%s\n```",
		tmpl,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			b.BuildCancelButton(),
		),
	)

	b.AddChatState(chatID, StateWaitingUpdateCompanyIntegration)
	b.AddPendingCompanyID(chatID, params.CompanyID)

	b.SendPendingMessage(b.NewKeyboardMessage(chatID, msgText, keyboard))
}

func (b *Bot) HandleUpdateCompanyIntegration(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	companyID := b.pendingCompanyIDs[chatID]

	company, err := b.CompanySvc.GetByID(companyID)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	_, err = b.CompanySvc.CreateOrUpdateCompanyIntegration(company.ID, message.Text)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	config, err := b.BuildCompanyIntegrationsDetail(chatID, company)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.ClearChatState(chatID)

	b.SendMessage(*config)
}
