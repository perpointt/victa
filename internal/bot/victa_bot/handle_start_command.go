package victa_bot

import (
	"context"
	"errors"
	_ "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"victa/internal/domain"

	appErr "victa/internal/errors"
)

func (b *Bot) HandleStartCommand(ctx context.Context, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	tgID := message.From.ID
	name := message.From.FirstName
	payload := message.CommandArguments()

	existing, err := b.UserSvc.GetByTgID(ctx, tgID)
	if err != nil && !errors.Is(err, appErr.ErrUserNotFound) {
		b.SendErrorMessage(chatID, err)
		return
	}

	var user *domain.User
	if existing == nil {
		user, err = b.UserSvc.Register(ctx, tgID, name)
		if err != nil {
			b.SendErrorMessage(chatID, err)
			return
		}
		existing = user
	} else {
		user, err = b.UserSvc.Update(ctx, existing.ID, name)
		if err != nil {
			b.SendErrorMessage(chatID, err)
			return
		}
		existing = user
	}

	if payload != "" {
		companyID, err := b.InviteSvc.ValidateToken(payload)
		if err == nil {
			if err := b.CompanySvc.AddUserToCompany(ctx, existing.ID, companyID); err != nil {
				b.SendErrorMessage(chatID, err)
			}
		} else {
			b.SendErrorMessage(chatID, err)
		}
	}

	menu, err := b.BuildMainMenu(ctx, chatID, user)
	if err != nil {
		b.SendErrorMessage(chatID, err)
		return
	}

	b.SendMessage(*menu)
}
