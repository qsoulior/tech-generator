package task_process_handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rabbitmq/amqp091-go"

	"github.com/qsoulior/tech-generator/backend/internal/usecase/task_process/domain"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Handle(ctx context.Context, msg amqp091.Delivery) error {
	body := string(msg.Body)

	taskID, err := strconv.ParseInt(body, 10, 64)
	if err != nil {
		_ = msg.Nack(false, false)
		return fmt.Errorf("strconv - parse int: %w", err)
	}

	in := domain.TaskProcessIn{
		TaskID: taskID,
	}

	err = h.usecase.Handle(ctx, in)
	if err != nil {
		_ = msg.Nack(false, false)
		return fmt.Errorf("task process usecase: %w", err)
	}

	_ = msg.Ack(false)
	return nil
}
