package bot

import (
	"fmt"
	_ "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

// HandleStartCommand обрабатывает /start, учитывая необязательный payload-invite
func (b *Bot) HandleStartCommand(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	tgID := message.From.ID
	name := message.From.FirstName
	payload := message.CommandArguments()

	existing, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		b.SendMessage(b.NewMessage(chatID, "Ошибка при проверке пользователя."))
		return
	}

	var user *domain.User
	if existing == nil {
		// новый пользователь
		user, err = b.UserSvc.Register(fmt.Sprintf("%d", tgID), name)
		if err != nil {
			b.SendMessage(b.NewMessage(chatID, "Ошибка регистрации пользователя."))
			return
		}
		existing = user
	} else {
		//// обновляем имя
		//user, err = b.UserSvc.Update(existing.ID, name)
		//if err != nil {
		//	b.SendMessage(b.NewMessage(chatID, "Ошибка обновления данных пользователя."))
		//	return
		//}
	}

	if payload != "" {
		companyID, err := b.InviteSvc.ValidateToken(payload)
		if err == nil {
			if err := b.CompanySvc.AddUserToCompany(existing.ID, companyID); err != nil {
				b.SendMessage(b.NewMessage(chatID, "Ошибка при добавлении пользователя в компанию."))
			}
		} else {
			b.SendMessage(b.NewMessage(chatID, "Ошибка при проверке приглашения"))
		}
	}

	msg := b.BuildMainMenu(chatID, existing)
	if msg == nil {
		return
	}
	b.SendMessage(*msg)
}
