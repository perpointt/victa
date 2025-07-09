package bot_common

import (
	"sync"
	"victa/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotFactory struct {
	mu    sync.Mutex
	cache map[string]*BaseBot
}

func NewBotFactory() *BotFactory {
	return &BotFactory{
		cache: make(map[string]*BaseBot),
	}
}

// GetBaseBot возвращает *BotAPI для token, создавая и кэшируя при первом запросе.
func (f *BotFactory) GetBaseBot(token string, logger logger.Logger) (*BaseBot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if api, ok := f.cache[token]; ok {
		return api, nil
	}
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	baseBot := NewBaseBot(api, logger)
	f.cache[token] = baseBot

	return baseBot, nil
}
