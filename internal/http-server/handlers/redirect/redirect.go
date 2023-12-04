package redirect

import (
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	Alias string `json:"alias" validate:"required"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type GetterURL interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, getterUrl GetterURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "andlers.redirect.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("required_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed decode request"))

			return
		}

		log.Info("request body  decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("failed to validate request body", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias, err := getterUrl.GetURL(req.Alias)
		if err != nil {
			log.Error("failed to get alias", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get alias"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
