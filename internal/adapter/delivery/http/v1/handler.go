package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/authena-ru/courses-organization/internal/app"
)

type coursesOrganizationHandler struct {
	app app.Application
}

func NewHandler(app app.Application, r chi.Router) http.Handler {
	return HandlerFromMux(coursesOrganizationHandler{
		app: app,
	}, r)
}
