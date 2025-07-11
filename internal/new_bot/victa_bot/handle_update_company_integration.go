package victa_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
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
		b.SendMessage(b.NewMessage(chatID, "Ошибка при получении интеграций."))
		return
	}

	tmpl, err := b.BuildIntegrationTemplate(ci)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при создании шаблона"))
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
	tgID := message.From.ID
	companyID := b.pendingCompanyIDs[chatID]

	user, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при поиске пользователя."))
		return
	}

	company, err := b.CompanySvc.GetByID(companyID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при поиске компании."))
		return
	}

	_, err = b.CompanySvc.CreateOrUpdateCompanyIntegration(user.ID, message.Text)
	if err != nil {
		log.Fatalf(err.Error())
		b.SendMessage(b.NewMessage(chatID, "Ошибка при обновлении интеграции."))
		return
	}

	config, err := b.BuildCompanyIntegrationsDetail(chatID, company)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при получении интеграций."))
		return
	}

	b.SendMessage(*config)
}
