package http

import (
	auth2 "github.com/authena-ru/courses-organization/internal/port/http/auth"
	v12 "github.com/authena-ru/courses-organization/internal/port/http/v1"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/pkg/logging"
)

func NewHandler(app app.Application) http.Handler {
	apiRouter := chi.NewRouter()
	addMiddlewares(apiRouter)

	rootRouter := chi.NewRouter()
	rootRouter.Mount("/v1", v12.NewHandler(app, apiRouter))

	return rootRouter
}

func addMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logging.NewStructuredLogger(logrus.StandardLogger()))
	router.Use(middleware.Recoverer)
	addCORSMiddleware(router)
	router.Use(auth2.MockAuthHTTPMiddleware)
	router.Use(middleware.NoCache)
}

const maxAge = 300

func addCORSMiddleware(router *chi.Mux) {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ";")
	if len(allowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           maxAge,
	})
	router.Use(corsMiddleware.Handler)
}
