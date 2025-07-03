package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (b *Bot) AddChatState(chatID int64, state ChatState) {
	states[chatID] = state
}

func (b *Bot) DeleteChatState(chatID int64) {
	delete(states, chatID)
}

func (b *Bot) SendPendingMessage(config tgbotapi.MessageConfig) {
	sentMsg, err := b.SendMessage(config)
	if err == nil {
		pendingMessages[config.ChatID] = sentMsg.MessageID
	}
}

func (b *Bot) DeletePendingMessage(chatID int64) {
	if msgID, ok := pendingMessages[chatID]; ok {
		err := b.DeleteMessage(chatID, msgID)
		if err == nil {
			delete(pendingMessages, chatID)
		}
	}
}

func (b *Bot) CancelAllOperation(chatID int64) {
	b.CancelCreateCompany(chatID)
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
