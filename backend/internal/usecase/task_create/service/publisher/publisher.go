package publisher

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rabbitmq/amqp091-go"
)

type Service struct {
	amqpPublisher amqpPublisher
}

func New(amqpPublisher amqpPublisher) *Service {
	return &Service{
		amqpPublisher: amqpPublisher,
	}
}

func (s *Service) PublishTaskCreated(ctx context.Context, id int64) error {
	body := strconv.FormatInt(id, 10)

	msg := amqp091.Publishing{
		DeliveryMode: amqp091.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	}

	err := s.amqpPublisher.PublishWithContext(ctx, "", "task_created", false, false, msg)
	if err != nil {
		return fmt.Errorf("amqp publisher - publish with context: %w", err)
	}

	return nil
}
