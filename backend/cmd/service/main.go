package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/qsoulior/tech-generator/backend/internal/pkg/httpserver"
	"github.com/qsoulior/tech-generator/backend/internal/pkg/postgres"
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
		logger.Error("close postgres connection", slog.String("err", err.Error()))
		code = 1
	}()

	server := httpserver.New(nil, logger)
	if err := server.Run(ctx); err != nil {
		logger.Error("fail server", slog.String("err", err.Error()))
		return 1
	}

	logger.Info("stop server")
	return 0
}
