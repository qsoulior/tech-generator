package error_handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen/ogenerrors"
)

type Handler struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	code := ogenerrors.ErrorCode(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if code == http.StatusInternalServerError {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("internal server error",
			slog.String("path", r.URL.Path),
			slog.String("err", err.Error()))
		return
	}

	e := jx.GetEncoder()
	e.ObjStart()
	e.FieldStart("error")
	e.StrEscape(err.Error())
	e.ObjEnd()

	_, _ = w.Write(e.Bytes())
	jx.PutEncoder(e)
}
