//go:generate go tool mockgen -package $GOPACKAGE -source contract.go -destination contract_mock.go

package publisher

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type amqpPublisher interface {
	PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error
}
