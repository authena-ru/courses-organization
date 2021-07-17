package v1_test

import (
	"bytes"
	auth2 "github.com/authena-ru/courses-organization/internal/port/http/auth"
	v12 "github.com/authena-ru/courses-organization/internal/port/http/v1"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
	"github.com/authena-ru/courses-organization/pkg/logging"
)

func newHTTPHandler(t *testing.T, application app.Application) http.Handler {
	t.Helper()

	router := chi.NewRouter()
	router.Use(logging.NewStructuredLogger(logrus.StandardLogger()))

	return v12.NewHandler(application, router)
}

func newHTTPRequest(t *testing.T, method, target, body string, authorized course.Academic) *http.Request {
	t.Helper()

	var r *http.Request
	if body != "" {
		r = httptest.NewRequest(method, target, bytes.NewBufferString(body))
		r.Header.Set("Content-Type", "application/json")
	} else {
		r = httptest.NewRequest(method, target, nil)
	}

	r = r.WithContext(auth2.WithAcademicInCtx(r.Context(), authorized))

	return r
}
