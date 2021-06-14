package http

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"

	"github.com/authena-ru/courses-organization/internal/coursesorg/adapter/delivery/http/auth"
	"github.com/authena-ru/courses-organization/internal/coursesorg/adapter/delivery/http/logging"
	"github.com/authena-ru/courses-organization/internal/coursesorg/adapter/delivery/http/v1"
	"github.com/authena-ru/courses-organization/internal/coursesorg/app"
)

func NewHandler(app app.Application) http.Handler {
	apiRouter := chi.NewRouter()
	addMiddlewares(apiRouter)

	rootRouter := chi.NewRouter()
	rootRouter.Mount("/v1", v1.NewHandler(app, apiRouter))
	return rootRouter
}

func addMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logging.NewStructuredLogger(logrus.StandardLogger()))
	router.Use(middleware.Recoverer)
	addCORSMiddleware(router)
	router.Use(auth.MockAuthHTTPMiddleware)
	router.Use(middleware.NoCache)
}

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
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)
}
