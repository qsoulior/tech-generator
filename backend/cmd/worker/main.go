package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/qsoulior/tech-generator/backend/internal/pkg/postgres"
	"github.com/qsoulior/tech-generator/backend/internal/pkg/rabbitmq"
	task_process_handler "github.com/qsoulior/tech-generator/backend/internal/transport/amqp/handler/task_process"
	task_process_usecase "github.com/qsoulior/tech-generator/backend/internal/usecase/task_process"
)

func main() {
	os.Exit(run())
}

func run() (code int) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	err := godotenv.Overload()
	if err != nil {
		logger.Error("overload env", slog.String("err", err.Error()))
		return 1
	}

	db, err := postgres.Connect(ctx)
	if err != nil {
		logger.Error("connect postgres", slog.String("err", err.Error()))
		return 1
	}
	defer func() {
		err := db.Close()
		if err != nil {
			logger.Error("close postgres connection", slog.String("err", err.Error()))
			code = 1
		}
	}()

	conn, err := rabbitmq.Connect()
	if err != nil {
		logger.Error("connect rabbitmq", slog.String("err", err.Error()))
		return 1
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Error("close rabbitmq connection", slog.String("err", err.Error()))
			code = 1
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("connect  rabbitmq channel", slog.String("err", err.Error()))
		return 1
	}
	defer func() {
		err := ch.Close()
		if err != nil {
			logger.Error("close rabbitmq channel", slog.String("err", err.Error()))
			code = 1
		}
	}()

	_, err = ch.QueueDeclare("task_created", true, false, false, false, nil)
	if err != nil {
		logger.Error("declare queue", slog.String("err", err.Error()))
		return 1
	}

	taskProcessUsecase := task_process_usecase.New(db)
	taskProcessHandler := task_process_handler.New(taskProcessUsecase)

	msgs, err := ch.ConsumeWithContext(ctx, "task_created", "", false, false, false, false, nil)
	if err != nil {
		logger.Error("consume task_created", slog.String("err", err.Error()))
		return 1
	}

	for msg := range msgs {
		err := taskProcessHandler.Handle(ctx, msg)
		if err != nil {
			logger.Error("task process handler", slog.String("err", err.Error()))
		}
	}

	logger.Info("stop consumer")
	return 0
}
