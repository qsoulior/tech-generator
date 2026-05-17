package auth_middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-faster/jx"

	error_domain "github.com/qsoulior/tech-generator/backend/internal/domain/error"
	"github.com/qsoulior/tech-generator/backend/internal/generated/api"
)

const headerUserID = "X-User-Id"

var nonAuth = map[string]struct{}{
	"userTokenCreate": {},
	"userCreate":      {},
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

func (m *Middleware) Handle(next *api.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del(headerUserID)

		route, ok := next.FindRoute(r.Method, r.URL.Path)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		if _, ok := nonAuth[route.OperationID()]; ok {
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
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)

				e := jx.GetEncoder()
				e.ObjStart()
				e.FieldStart("error")
				e.StrEscape(baseErr.Error())
				e.ObjEnd()
				_, _ = w.Write(e.Bytes())
				jx.PutEncoder(e)
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
