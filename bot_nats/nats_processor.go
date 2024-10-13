package bot_nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

func PublishToNATS(natsQueue string, message []byte) error {
	nc, err := Connect()
	if err != nil {
		return err
	}
	defer nc.Close()

	err = nc.Publish(natsQueue, message)
	if err != nil {
		return err
	}

	log.Printf("[PublishToNATS] Повідомлення надіслано в NATS")
	return nil
}
func Connect() (*nats.Conn, error) {
	ip := os.Getenv("BROKER_IP")
	port := os.Getenv("BROKER_PORT")
	natsUrl := fmt.Sprintf("nats://%s:%s", ip, port)

	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, err
	}
	return nc, nil
}
