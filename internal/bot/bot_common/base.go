package bot_common

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type BaseBot struct {
	BotAPI *tgbotapi.BotAPI
}

func NewBaseBot(api *tgbotapi.BotAPI) *BaseBot {
	return &BaseBot{BotAPI: api}
}

func (b *BaseBot) AnswerCallback(callback *tgbotapi.CallbackQuery, text string) {
	answer := tgbotapi.NewCallback(callback.ID, text)
	if _, err := b.BotAPI.Request(answer); err != nil {
		log.Printf(err.Error())
	}
}

func (b *BaseBot) SendErrorMessage(config tgbotapi.MessageConfig) *tgbotapi.Message {
	log.Printf(config.Text)

	msg, err := b.BotAPI.Send(config)

	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
		return nil
	}

	return &msg
}

func (b *BaseBot) SendMessage(config tgbotapi.MessageConfig) *tgbotapi.Message {
	msg, err := b.BotAPI.Send(config)

	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
		return nil
	}

	return &msg
}

func (b *BaseBot) EditMessage(messageID int, config tgbotapi.MessageConfig) *tgbotapi.Message {
	text := config.Text
	replyMarkup, ok := config.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)

	if !ok {
		return nil
	}

	editMsg := tgbotapi.NewEditMessageTextAndMarkup(config.ChatID, messageID, text, replyMarkup)
	editMsg.ParseMode = config.ParseMode

	msg, err := b.BotAPI.Send(editMsg)

	if err != nil {
		b.SendErrorMessage(b.NewMessage(config.ChatID, err.Error()))
		return nil
	}

	return &msg
}

func (b *BaseBot) DeleteMessage(chatID int64, messageID int) {
	delMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	if _, err := b.BotAPI.Request(delMsg); err != nil {
		b.SendErrorMessage(b.NewMessage(chatID, err.Error()))
	}
}

func (b *BaseBot) NewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

func (b *BaseBot) NewHtmlMessage(chatID int64, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true
	return msg
}

func (b *BaseBot) NewKeyboardMessage(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboard
	return msg
}
