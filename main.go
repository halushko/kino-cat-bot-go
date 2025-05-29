package main

import (
	"github.com/halushko/kino-cat-core-go/logger_helper"
	"github.com/halushko/kino-cat-core-go/nats_helper"
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

	basePoller := &telebot.LongPoller{Timeout: 1 * time.Second}
	mwPoller := telebot.NewMiddlewarePoller(
		basePoller,
		func(upd *telebot.Update) bool {
			urls := ""
			if upd.Message != nil && len(upd.Message.CaptionEntities) > 0 {
				for _, entity := range upd.Message.CaptionEntities {
					if entity.URL != "" {
						urls = urls + entity.URL + "\n"
					}
				}
			}
			if urls != "" {
				nats_helper.PublishTextMessage("TELEGRAM_INPUT_TEXT_QUEUE", upd.Message.Chat.ID, urls)
				return false
			}
			return true
		})

	pref := telebot.Settings{
		Token:  token,
		Poller: mwPoller,
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
