package handlers

import (
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"gopkg.in/telebot.v3"
	"log"
)

func HandleTextMessages(bot *telebot.Bot) {
	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		chatId := c.Chat().ID
		message := c.Message().Text

		log.Printf("[HandleTextMessages] chatId:%d, message:%s", chatId, message)

		nats_helper.PublishTextMessage("TELEGRAM_INPUT_TEXT_QUEUE", chatId, message)

		if err := c.Send("Ваше повідомлення " + message + " додано до обробки"); err != nil {
			log.Printf("[HandleTextMessages] ERROR:%v", err)
			return err
		}

		return nil
	})
}
