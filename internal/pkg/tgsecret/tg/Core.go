package tg

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	*tgbotapi.BotAPI
	ChatNotificationID int64
}

// Создать приложение телеграмма
func NewTelegram(cf Config) (*Telegram, error) {

	bot, Err := tgbotapi.NewBotAPI(cf.Token)
	if Err != nil {
		return nil, fmt.Errorf("os.ReadFile: %v", Err)
	}
	bot.Debug = false

	return &Telegram{BotAPI: bot, ChatNotificationID: cf.ChatID}, nil
}
