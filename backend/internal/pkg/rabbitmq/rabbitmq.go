package rabbitmq

import (
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp091.Connection, error) {
	url := os.Getenv("RABBITMQ_URL")

	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
