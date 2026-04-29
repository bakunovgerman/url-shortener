package delete

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shorter/base/internal/lib/api/response"
	"url-shorter/base/internal/lib/logger/sl"
	"url-shorter/base/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) (string, error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		_, err := urlDeleter.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.AliasNotFound) {
				log.Error("alias not found", sl.Err(err))
				render.JSON(w, r, resp.Error("alias not found"))
				return
			}

			log.Error("failed to delete", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete"))
			return
		}

		log.Info("alias deleted", slog.String("alias", alias))
		render.JSON(w, r, Response{Response: resp.OK(), Alias: alias})
	}
}
