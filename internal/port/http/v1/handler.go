package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/authena-ru/courses-organization/internal/app"
)

type handler struct {
	app app.Application
}

func NewHandler(app app.Application, r chi.Router) http.Handler {
	return HandlerFromMux(handler{
		app: app,
	}, r)
}
