package v1_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
	"github.com/authena-ru/courses-organization/internal/port/http/auth"
	v1 "github.com/authena-ru/courses-organization/internal/port/http/v1"
	"github.com/authena-ru/courses-organization/pkg/logging"
)

func newHTTPHandler(t *testing.T, application app.Application) http.Handler {
	t.Helper()

	router := chi.NewRouter()
	router.Use(logging.NewStructuredLogger(logrus.StandardLogger()))

	return v1.NewHandler(application, router)
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

	r = r.WithContext(auth.WithAcademicInCtx(r.Context(), authorized))

	return r
}

func requireAddTaskCommandsEquals(t *testing.T, expectedCommand, givenCommand app.AddTaskCommand) {
	t.Helper()

	require.Equalf(t, expectedCommand.Academic, givenCommand.Academic, "academics are not equal")
	require.Equalf(t, expectedCommand.CourseID, givenCommand.CourseID, "course ID's are not equal")
	require.Equalf(t, expectedCommand.TaskTitle, givenCommand.TaskTitle, "task titles are not equal")
	require.Equalf(t, expectedCommand.TaskDescription, givenCommand.TaskDescription, "task descriptions are not equal")
	require.Equalf(t, expectedCommand.TaskType, givenCommand.TaskType, "task types are not equal")
	requireDatesEquals(t, expectedCommand.Deadline.GoodGradeTime(), givenCommand.Deadline.GoodGradeTime())
	requireDatesEquals(t, expectedCommand.Deadline.ExcellentGradeTime(), givenCommand.Deadline.ExcellentGradeTime())
}

func requireDatesEquals(t *testing.T, expectedDate, givenDate time.Time) {
	t.Helper()

	expectedYear, expectedMonth, expectedDay := expectedDate.Date()
	givenYear, givenMonth, givenDay := givenDate.Date()

	require.Equalf(t, expectedYear, givenYear, "years are not equal")
	require.Equalf(t, expectedMonth, givenMonth, "months are not equal")
	require.Equalf(t, expectedDay, givenDay, "days are not equal")
}
