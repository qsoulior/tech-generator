package task_process_handler

import (
	"context"
	"errors"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

type ackCall struct {
	tag      uint64
	multiple bool
}

type nackCall struct {
	tag      uint64
	multiple bool
	requeue  bool
}

type fakeAcknowledger struct {
	acks  []ackCall
	nacks []nackCall
}

func (a *fakeAcknowledger) Ack(tag uint64, multiple bool) error {
	a.acks = append(a.acks, ackCall{tag: tag, multiple: multiple})
	return nil
}

func (a *fakeAcknowledger) Nack(tag uint64, multiple, requeue bool) error {
	a.nacks = append(a.nacks, nackCall{tag: tag, multiple: multiple, requeue: requeue})
	return nil
}

func (a *fakeAcknowledger) Reject(tag uint64, requeue bool) error {
	return nil
}

func newDelivery(body string) (amqp091.Delivery, *fakeAcknowledger) {
	ack := &fakeAcknowledger{}
	msg := amqp091.Delivery{
		Acknowledger: ack,
		DeliveryTag:  42,
		Body:         []byte(body),
	}
	return msg, ack
}

func TestHandler_Handle_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewMockusecase(ctrl)
	usecase.EXPECT().Handle(ctx, domain.TaskProcessIn{TaskID: 1234}).Return(nil)

	msg, ack := newDelivery("1234")

	handler := New(usecase)
	err := handler.Handle(ctx, msg)
	require.NoError(t, err)

	require.Len(t, ack.acks, 1)
	require.Equal(t, ackCall{tag: 42, multiple: false}, ack.acks[0])
	require.Empty(t, ack.nacks)
}

func TestHandler_Handle_Error(t *testing.T) {
	ctx := context.Background()

	testErr := errors.New("test error")

	tests := []struct {
		name     string
		body     string
		setup    func(usecase *Mockusecase)
		wantErr  string
		wantAck  bool
		wantNack bool
	}{
		{
			name:     "strconv_ParseInt",
			body:     "not-a-number",
			setup:    func(usecase *Mockusecase) {},
			wantErr:  "strconv - parse int",
			wantAck:  false,
			wantNack: true,
		},
		{
			name: "usecase_Handle",
			body: "1234",
			setup: func(usecase *Mockusecase) {
				usecase.EXPECT().Handle(ctx, domain.TaskProcessIn{TaskID: 1234}).Return(testErr)
			},
			wantErr:  "task process usecase",
			wantAck:  false,
			wantNack: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usecase := NewMockusecase(ctrl)
			tt.setup(usecase)

			msg, ack := newDelivery(tt.body)

			handler := New(usecase)
			err := handler.Handle(ctx, msg)
			require.ErrorContains(t, err, tt.wantErr)

			if tt.wantAck {
				require.Len(t, ack.acks, 1)
			} else {
				require.Empty(t, ack.acks)
			}
			if tt.wantNack {
				require.Len(t, ack.nacks, 1)
				require.Equal(t, nackCall{tag: 42, multiple: false, requeue: false}, ack.nacks[0])
			} else {
				require.Empty(t, ack.nacks)
			}
		})
	}
}
