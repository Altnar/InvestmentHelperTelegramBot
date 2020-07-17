package loggerbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type LoggerBot struct {
	LoggerBotAPI *tgbotapi.BotAPI
	ChatID       int64
}

func NewLoggerBot(botAPI string, chatID int) LoggerBot {
	logBot, _ := tgbotapi.NewBotAPI(botAPI)
	return LoggerBot{logBot, int64(chatID)}
}

func (bot *LoggerBot) SendLog(logText string) error {
	logMessage := tgbotapi.NewMessage(bot.ChatID, logText)
	_, err := bot.LoggerBotAPI.Send(logMessage)
	if err != nil {
		return err
	}
	return nil
}
