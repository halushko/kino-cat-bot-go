package main

import (
	"github.com/halushko/kino-cat-core-go/logger_helper"
	"kino-cat-bot-go/handlers"
	"kino-cat-bot-go/listeners"
	"log"
	"os"
	"time"

	"gopkg.in/telebot.v3"
)

func main() {
	logFile := logger_helper.SoftPrepareLogFile()
	bot := prepareBot()
	listeners.StartTextMessagesSender(bot)
	log.Println("Бота запущено")
	bot.Start()
	logger_helper.SoftLogClose(logFile)
}

func prepareBot() *telebot.Bot {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("Необхідно задати токен боту в env BOT_TOKEN")
	}
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 1 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	registerBotHandlers(bot)

	return bot
}

func registerBotHandlers(bot *telebot.Bot) {
	handlers.HandleTextMessages(bot)
	handlers.HandleDocuments(bot)
}
