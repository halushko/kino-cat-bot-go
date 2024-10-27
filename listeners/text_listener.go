package listeners

import (
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"gopkg.in/telebot.v3"
	"log"
)

func StartTextMessagesSender(bot *telebot.Bot) {

	processor := func(data []byte) {
		log.Printf("[StartTextMessagesSender] Отримано повідомлення з NATS: %s", string(data))
		userId, messageText, err := nats_helper.ParseNatsBotText(data)

		if err != nil {
			log.Printf("[StartTextMessagesSender] Помилка при розборі повідомлення з NATS: %v", err)
			return
		}

		log.Printf("[StartTextMessagesSender] Парсинг повідомлення: chatID = %d, message = %s", userId, messageText) // Новый лог для проверки данных

		if userId != 0 && messageText != "" {
			_, err := bot.Send(&telebot.User{ID: userId}, messageText)
			if err != nil {
				log.Printf("[StartTextMessagesSender] Помилка при відправленні повідомлення користувачу: %v", err)
			} else {
				log.Printf("[StartTextMessagesSender] Повідомлення надіслане користовачу: chatID = %d, message = %s", userId, messageText)
			}
		} else {
			log.Printf("[StartTextMessagesSender] Помилка: ID користувача чи текст повідомлення порожні")
		}
	}

	listener := &nats_helper.NatsListenerHandler{
		Function: processor,
	}

	err := nats_helper.StartNatsListener("TELEGRAM_OUTPUT_TEXT_QUEUE", listener)
	if err != nil {
		log.Printf("[StartTextMessagesSender] Помилка: %v", err)
	}
}
