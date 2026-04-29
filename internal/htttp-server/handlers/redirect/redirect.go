package redirect

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

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, resp.Error("not found"))
			return
		}

		if err != nil {
			log.Error("failed to get url", "alias", alias, sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}
		log.Info("got url", "alias", alias, slog.String("url", resURL))
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
