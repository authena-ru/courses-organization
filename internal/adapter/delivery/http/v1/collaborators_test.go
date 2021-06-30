package v1

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/auth"
	"github.com/authena-ru/courses-organization/internal/adapter/delivery/http/logging"
	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/app/command/mock"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestHandler_AddCollaboratorToCourse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		RequestBody    string
		Authorized     course.Academic
		CourseID       string
		Command        app.AddCollaboratorCommand
		PrepareHandler func(expectedCommand app.AddCollaboratorCommand) mock.AddCollaboratorsHandler
		StatusCode     int
		ResponseBody   string
	}{
		{
			Name:        "collaborator_added_to_course",
			RequestBody: `{"id": "199cf094-0b92-455a-9da3-f353f4bf9ed3"}`,
			Authorized:  course.MustNewAcademic("1009d0ed-600f-4bd1-96fa-8ccaedb4e7d7", course.TeacherType),
			CourseID:    "ecea7dcc-a1d9-48cc-8526-0d4b58bc298b",
			Command: app.AddCollaboratorCommand{
				CourseID:       "ecea7dcc-a1d9-48cc-8526-0d4b58bc298b",
				Academic:       course.MustNewAcademic("1009d0ed-600f-4bd1-96fa-8ccaedb4e7d7", course.TeacherType),
				CollaboratorID: "199cf094-0b92-455a-9da3-f353f4bf9ed3",
			},
			PrepareHandler: func(expectedCommand app.AddCollaboratorCommand) mock.AddCollaboratorsHandler {
				return func(_ context.Context, givenCommand app.AddCollaboratorCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return nil
				}
			},
			StatusCode:   http.StatusNoContent,
			ResponseBody: "",
		},
		{
			Name:        "bad_request",
			RequestBody: `"id": ""`,
			Authorized:  course.MustNewAcademic("1ccb6e85-80ed-4f4d-aa76-5910ad054820", course.TeacherType),
			CourseID:    "556a7670-e0cc-4867-af4d-b3d142bc0f56",
			Command: app.AddCollaboratorCommand{
				CourseID:       "556a7670-e0cc-4867-af4d-b3d142bc0f56",
				Academic:       course.MustNewAcademic("1ccb6e85-80ed-4f4d-aa76-5910ad054820", course.TeacherType),
				CollaboratorID: "",
			},
			PrepareHandler: func(expectedCommand app.AddCollaboratorCommand) mock.AddCollaboratorsHandler {
				return func(_ context.Context, givenCommand app.AddCollaboratorCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return nil
				}
			},
			StatusCode: http.StatusBadRequest,
			ResponseBody: `{
								"slug": "bad-request",
								"details": "json: cannot unmarshal string into Go value of type v1.AddCollaboratorToCourseRequest"
							}`,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			router := chi.NewRouter()
			router.Use(logging.NewStructuredLogger(logrus.StandardLogger()))

			application := app.Application{
				Commands: app.Commands{
					AddCollaborator: c.PrepareHandler(c.Command),
				},
			}
			h := NewHandler(application, router)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPut,
				fmt.Sprintf("/courses/%s/collaborators", c.CourseID),
				bytes.NewBufferString(c.RequestBody),
			)
			r = r.WithContext(auth.WithAcademicInCtx(r.Context(), c.Authorized))
			r.Header.Set("Content-Type", "application/json")

			h.ServeHTTP(w, r)

			require.Equal(t, c.StatusCode, w.Code)

			if c.ResponseBody != "" {
				require.JSONEq(t, c.ResponseBody, w.Body.String())
			}
		})
	}
}
