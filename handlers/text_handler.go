package handlers

import (
	"encoding/json"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"gopkg.in/telebot.v3"
	"log"
)

type telegramMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func HandleTextMessages(bot *telebot.Bot) {
	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		chatId := c.Chat().ID
		message := c.Message().Text

		log.Printf("[HandleTextMessages] chatId:%d, message:%s", chatId, message)

		msg := telegramMessage{
			ChatID: chatId,
			Text:   message,
		}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			log.Printf("[HandleTextMessages] ERROR:%s", err)
			return err
		}

		if err = nats_helper.PublishToNATS("TELEGRAM_INPUT_TEXT_QUEUE", jsonData); err != nil {
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
