package notification_bot

import (
	"strconv"
	"victa/internal/new_bot/bot_common"
)

// Bot хранит API и ссылку на БД
type Bot struct {
	*bot_common.BaseBot
	chatID int64
}

// NewBot создаёт нового бота
func NewBot(
	base *bot_common.BaseBot,
	chatIDString string,
) (*Bot, error) {
	chatID, err := strconv.ParseInt(chatIDString, 10, 64)
	if err != nil {
		return nil, err
	}

	return &Bot{
		BaseBot: base,
		chatID:  chatID,
	}, nil
}
