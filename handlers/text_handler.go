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
			log.Printf("[HandleTextMessages] ERROR:%s", err)
			return err
		}

		if err = bot_nats.PublishToNATS("TELEGRAM_INPUT_TEXT_QUEUE", jsonData); err != nil {
			log.Printf("[HandleTextMessages] ERROR:%s", err)
			return err
		}

		if err = c.Send("Ваше повідомлення " + message + " додано до обробки"); err != nil {
			log.Printf("[HandleTextMessages] ERROR:%s", err)
			return err
		}

		return nil
	})
}
