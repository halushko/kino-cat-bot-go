package main

import (
	"kino-cat-bot-go/handlers"
	"log"
	"os"
	"time"

	"gopkg.in/telebot.v3"
)

func main() {
	log.SetOutput(prepareLogFile())
	bot := prepareBot()

	log.Println("Бот запущен")
	bot.Start()
}

func prepareLogFile() *os.File {
	log.Print("Старт бота")

	logFile, err := os.OpenFile("bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Printf("Помилка при спробі закрити лог файл: %v", err)
		}
	}()

	return logFile
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
