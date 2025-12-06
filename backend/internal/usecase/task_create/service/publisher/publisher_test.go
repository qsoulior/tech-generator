package publisher

import (
	"context"
	"errors"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_PublishTaskCreated_Success(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	amqpPublisher := NewMockamqpPublisher(ctrl)

	id := int64(1234)
	msg := amqp091.Publishing{
		DeliveryMode: amqp091.Persistent,
		ContentType:  "text/plain",
		Body:         []byte("1234"),
	}

	amqpPublisher.EXPECT().PublishWithContext(ctx, "", "task_created", false, false, msg).Return(nil)

	service := New(amqpPublisher)
	err := service.PublishTaskCreated(ctx, id)
	require.NoError(t, err)
}

func TestService_PublishTaskCreated_Error(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	amqpPublisher := NewMockamqpPublisher(ctrl)

	id := int64(1234)
	expectedErr := errors.New("test")

	amqpPublisher.EXPECT().PublishWithContext(ctx, "", "task_created", false, false, gomock.Any()).Return(expectedErr)

	service := New(amqpPublisher)
	err := service.PublishTaskCreated(ctx, id)
	require.ErrorIs(t, err, expectedErr)
}
