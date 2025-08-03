package httpserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/qsoulior/tech-generator/backend/internal/pkg/env"
)

type Server struct {
	server *http.Server
	logger *slog.Logger
}

func New(handler http.Handler, logger *slog.Logger, opts ...OptionFunc) *Server {
	host := env.GetString("SERVICE_HOST", "localhost")
	port := env.GetString("SERVICE_PORT", "3000")

	httpServer := &http.Server{
		Handler:      handler,
		Addr:         net.JoinHostPort(host, port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	server := &Server{
		server: httpServer,
		logger: logger,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *Server) Run(ctx context.Context) error {
	s.server.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- s.server.ListenAndServe()
		close(errCh)
	}()

	s.logger.Info("start server", slog.String("addr", s.server.Addr))

	select {
	case <-ctx.Done():
		return s.server.Shutdown(ctx)
	case err := <-errCh:
		return err
	}
}
