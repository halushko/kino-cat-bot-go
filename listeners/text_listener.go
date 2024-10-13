package listeners

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"gopkg.in/telebot.v3"
	"kino-cat-bot-go/bot_nats"
	"log"
)

func StartNatsListener(bot *telebot.Bot) *nats.Conn {
	nc, err := bot_nats.Connect()
	if err != nil {
		log.Printf("Помилка при підключенні до NATS:", err)
	}
	_, err = nc.Subscribe("TELEGRAM_OUTPUT_TEXT_QUEUE", func(msg *nats.Msg) {
		log.Println("Отримано повідомлення з NATS: %s", string(msg.Data))
		chatID, messageText := parseNatsMessage(msg.Data)

		log.Println("Парсинг повідомлення: chatID = %d, message = %s", chatID, messageText) // Новый лог для проверки данных

		if chatID != 0 && messageText != "" {
			_, err := bot.Send(&telebot.User{ID: chatID}, messageText)
			if err != nil {
				log.Println("Помилка при відправленні повідомлення користувачу: %v", err)
			} else {
				log.Printf("Повідомлення надіслане користовачу: chatID = %d, message = %s", chatID, messageText)
			}
		} else {
			log.Println("Помилка: ID користувача чи текст повідомлення порожні")
		}
	})

	if err != nil {
		log.Println("Помилка підписки до черги NATS:", err)
	}

	err = nc.Flush()
	if err != nil {
		log.Println("Помилка після підписки до черги NATS:", err)
		return nil
	}
	if err = nc.LastError(); err != nil {
		log.Println("Помилка після підписки до черги NATS:", err)
	}

	log.Println("Підписка до черги NATS виконана")
	return nc
}

func parseNatsMessage(data []byte) (int64, string) {
	type NatsMessage struct {
		ChatID int64  `json:"chat_id"`
		Text   string `json:"text"`
	}

	var msg NatsMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("Помилка при розборы повыдомлення з NATS: %v", err)
		return 0, ""
	}

	return msg.ChatID, msg.Text
}
