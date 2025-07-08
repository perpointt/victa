package bot_common

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type BaseBot struct {
	BotAPI *tgbotapi.BotAPI
}

func NewBaseBot(api *tgbotapi.BotAPI) *BaseBot {
	return &BaseBot{BotAPI: api}
}

func (b *BaseBot) AnswerCallback(callback *tgbotapi.CallbackQuery, text string) error {
	answer := tgbotapi.NewCallback(callback.ID, text)
	if _, err := b.BotAPI.Request(answer); err != nil {
		log.Printf("Ошибка ответа на callback: %v", err)
		return err
	}

	return nil
}

func (b *BaseBot) SendMessage(config tgbotapi.MessageConfig) (*tgbotapi.Message, error) {
	msg, err := b.BotAPI.Send(config)

	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
		return nil, err
	}

	return &msg, nil
}

func (b *BaseBot) EditMessage(messageID int, config tgbotapi.MessageConfig) (*tgbotapi.Message, error) {
	text := config.Text
	replyMarkup, ok := config.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)

	if !ok {
		return nil, fmt.Errorf("expected InlineKeyboardMarkup, got %T", config.ReplyMarkup)
	}

	editMsg := tgbotapi.NewEditMessageTextAndMarkup(config.ChatID, messageID, text, replyMarkup)
	editMsg.ParseMode = config.ParseMode

	msg, err := b.BotAPI.Send(editMsg)

	if err != nil {
		log.Printf("Ошибка редактирования сообщения: %v", err)
		return nil, err
	}

	return &msg, nil
}

func (b *BaseBot) DeleteMessage(chatID int64, messageID int) error {
	delMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	if _, err := b.BotAPI.Request(delMsg); err != nil {
		log.Printf("Ошибка удаления сообщения: %v", err)
		return err
	}

	return nil
}

func (b *BaseBot) NewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

func (b *BaseBot) NewKeyboardMessage(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboard
	return msg
}
