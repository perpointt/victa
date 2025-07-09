package victa_bot

import (
	"fmt"
	_ "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"
)

func (b *Bot) HandleStartCommand(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	tgID := message.From.ID
	name := message.From.FirstName
	payload := message.CommandArguments()

	existing, err := b.UserSvc.GetByTgID(tgID)
	if err != nil {
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		return
	}

	var user *domain.User
	if existing == nil {
		user, err = b.UserSvc.Register(fmt.Sprintf("%d", tgID), name)
		if err != nil {
			b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
			return
		}
		existing = user
	} else {
		user, err = b.UserSvc.Update(existing.ID, name)
		if err != nil {
			b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
			return
		}
		existing = user
	}

	if payload != "" {
		companyID, err := b.InviteSvc.ValidateToken(payload)
		if err == nil {
			if err := b.CompanySvc.AddUserToCompany(existing.ID, companyID); err != nil {
				b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
			}
		} else {
			b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		}
	}

	menu, err := b.BuildMainMenu(chatID, user)
	if err != nil {
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
		return
	}

	b.SendMessage(*menu)
}
