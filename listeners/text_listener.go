package listeners

import (
	"encoding/json"
	"github.com/halushko/kino-cat-core-go"
	"github.com/nats-io/nats.go"
	"gopkg.in/telebot.v3"
	"log"
)

func StartTextMessagesSender(bot *telebot.Bot) {
	processor := func(msg *nats.Msg) {
		log.Println("[StartNatsListener] Отримано повідомлення з NATS: %s", string(msg.Data))
		chatID, messageText := parseNatsMessage(msg.Data)

		log.Println("[StartNatsListener] Парсинг повідомлення: chatID = %d, message = %s", chatID, messageText) // Новый лог для проверки данных

		if chatID != 0 && messageText != "" {
			_, err := bot.Send(&telebot.User{ID: chatID}, messageText)
			if err != nil {
				log.Println("[StartNatsListener] Помилка при відправленні повідомлення користувачу: %v", err)
			} else {
				log.Printf("[StartNatsListener] Повідомлення надіслане користовачу: chatID = %d, message = %s", chatID, messageText)
			}
		} else {
			log.Println("[StartNatsListener] Помилка: ID користувача чи текст повідомлення порожні")
		}
	}

	listener := &kino_cat_core_go.NatsListener{
		Handler: processor,
	}

	kino_cat_core_go.StartNatsListener("TELEGRAM_OUTPUT_TEXT_QUEUE", listener)
}

func parseNatsMessage(data []byte) (int64, string) {
	type NatsMessage struct {
		ChatID int64  `json:"chat_id"`
		Text   string `json:"text"`
	}

	var msg NatsMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("[StartNatsListener] Помилка при розборі повідомлення з NATS: %v", err)
		return 0, ""
	}

	return msg.ChatID, msg.Text
}
