package bot

import (
	"fmt"
	_ "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"victa/internal/domain"
)

// HandleStart обрабатывает /start, регистрирует или обновляет пользователя и сразу шлёт ответ
func (b *Bot) HandleStart(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	tgID := message.From.ID
	name := message.From.FirstName

	existing, err := b.UserSvc.FindByTgID(tgID)
	if err != nil {
		b.send(b.newMessage(chatID, "Ошибка при проверке пользователя."))
		log.Fatalf("Ошибка при проверке  пользователя: %v", err)
		return
	}

	var user *domain.User
	if existing == nil {
		user, err = b.UserSvc.Register(fmt.Sprintf("%d", tgID), name)
		if err != nil {
			b.send(b.newMessage(chatID, "Ошибка регистрации пользователя."))
			log.Fatalf("Ошибка регистрации пользователя: %v", err)
			return
		}

		existing = user
	}

	msg := b.BuildMainMenu(chatID, existing)
	b.send(msg)
}
