package handlers

import (
	"encoding/json"
	"gopkg.in/telebot.v3"
	"kino-cat-bot-go/bot_nats"
	"log"
)

type TelegramMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func HandleTextMessages(bot *telebot.Bot) {
	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		chatId := c.Chat().ID
		message := c.Message().Text

		log.Printf("[HandleTextMessages] chatId:%d, message:%s", chatId, message)

		msg := TelegramMessage{
			ChatID: chatId,
			Text:   message,
		}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		err = bot_nats.PublishToNATS("TELEGRAM_INPUT_TEXT_QUEUE", jsonData)
		if err != nil {
			return err
		}

		err = bot_nats.PublishToNATS("TELEGRAM_OUTPUT_TEXT_QUEUE", jsonData)
		if err != nil {
			log.Printf("[HandleTextMessages] Помилка при відправці повідомлення на TELEGRAM_OUTPUT_TEXT_QUEUE: %v", err)
			return err
		}
		log.Println("[HandleTextMessages] Повідомлення відправлено до TELEGRAM_OUTPUT_TEXT_QUEUE")
		return nil
	})
}
