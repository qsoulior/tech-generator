package auth_middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
)

const headerUserID = "X-User-Id"

var nonAuth = map[string]struct{}{
	"GET /user/create":        {},
	"POST /user/token/create": {},
}

type Middleware struct {
	usecase usecase
	logger  *slog.Logger
}

func New(usecase usecase, logger *slog.Logger) *Middleware {
	return &Middleware{
		usecase: usecase,
		logger:  logger,
	}
}

func (m *Middleware) Handle() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
			if _, ok := nonAuth[key]; ok {
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie("token")
			if err != nil || cookie.Valid() != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			user, err := m.usecase.Handle(r.Context(), cookie.Value)
			if err != nil {
				var baseErr *error_domain.BaseError
				if errors.As(err, &baseErr) {
					w.WriteHeader(http.StatusForbidden)
					_, _ = w.Write([]byte(err.Error()))
					return
				}

				m.logger.Error("user token parse", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Header.Set(headerUserID, strconv.FormatInt(user.ID, 10))
			next.ServeHTTP(w, r)
		})
	}
}
