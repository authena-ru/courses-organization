package v1_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/auth"
	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/logging"
	v1 "github.com/authena-ru/courses-organization/internal/adapter/delivery/http/v1"
	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func newHandler(t *testing.T, application app.Application) http.Handler {
	t.Helper()

	router := chi.NewRouter()
	router.Use(logging.NewStructuredLogger(logrus.StandardLogger()))

	return v1.NewHandler(application, router)
}

func newRequest(t *testing.T, method, target, requestBody string, authorized course.Academic) *http.Request {
	t.Helper()

	r := httptest.NewRequest(method, target, bytes.NewBufferString(requestBody))
	r = r.WithContext(auth.WithAcademicInCtx(r.Context(), authorized))
	r.Header.Set("Content-Type", "application/json")

	return r
}
