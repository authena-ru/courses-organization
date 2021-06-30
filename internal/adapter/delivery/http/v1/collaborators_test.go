package v1_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

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
		{
			Name:        "course_doesnt_exist",
			RequestBody: `{"id": "6db33767-f116-4499-89b2-3ef26fe842e3"}`,
			Authorized:  course.MustNewAcademic("69d13ada-be30-4c99-a93c-08cf1bce7eb8", course.TeacherType),
			CourseID:    "b59cc92a-574d-4065-87d1-955709b6964d",
			Command: app.AddCollaboratorCommand{
				CourseID:       "b59cc92a-574d-4065-87d1-955709b6964d",
				Academic:       course.MustNewAcademic("69d13ada-be30-4c99-a93c-08cf1bce7eb8", course.TeacherType),
				CollaboratorID: "6db33767-f116-4499-89b2-3ef26fe842e3",
			},
			PrepareHandler: func(expectedCommand app.AddCollaboratorCommand) mock.AddCollaboratorsHandler {
				return func(_ context.Context, givenCommand app.AddCollaboratorCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return app.ErrCourseDoesntExist
				}
			},
			StatusCode:   http.StatusNotFound,
			ResponseBody: `{"slug": "course-not-found", "details": "course doesn't exist"}`,
		},
		{
			Name:        "teacher_not_found",
			RequestBody: `{"id": "d825c9e0-abca-48f2-90a9-c53ef9636bde"}`,
			Authorized:  course.MustNewAcademic("aac3880c-46f2-44c9-9d2f-06e016124e48", course.TeacherType),
			CourseID:    "28104db1-8476-4279-830d-c49a6643a4b5",
			Command: app.AddCollaboratorCommand{
				CourseID:       "28104db1-8476-4279-830d-c49a6643a4b5",
				Academic:       course.MustNewAcademic("aac3880c-46f2-44c9-9d2f-06e016124e48", course.TeacherType),
				CollaboratorID: "d825c9e0-abca-48f2-90a9-c53ef9636bde",
			},
			PrepareHandler: func(expectedCommand app.AddCollaboratorCommand) mock.AddCollaboratorsHandler {
				return func(_ context.Context, givenCommand app.AddCollaboratorCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return app.ErrTeacherDoesntExist
				}
			},
			StatusCode:   http.StatusUnprocessableEntity,
			ResponseBody: `{"slug": "teacher-not-found", "details": "teacher doesn't exist"}`,
		},
		{
			Name:        "academic_cant_edit_course",
			RequestBody: `{"id": "22c381b2-d2d7-487b-b4f8-1f6b78f9cb92"}`,
			Authorized:  course.MustNewAcademic("d5cc070d-c562-4bcb-a6e8-e2af858d4a68", course.TeacherType),
			CourseID:    "b55d1633-f6d5-40a0-9bce-7a79f86518e5",
			Command: app.AddCollaboratorCommand{
				CourseID:       "b55d1633-f6d5-40a0-9bce-7a79f86518e5",
				Academic:       course.MustNewAcademic("d5cc070d-c562-4bcb-a6e8-e2af858d4a68", course.TeacherType),
				CollaboratorID: "22c381b2-d2d7-487b-b4f8-1f6b78f9cb92",
			},
			PrepareHandler: func(expectedCommand app.AddCollaboratorCommand) mock.AddCollaboratorsHandler {
				return func(_ context.Context, givenCommand app.AddCollaboratorCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return course.AcademicCantEditCourseError{}
				}
			},
			StatusCode:   http.StatusForbidden,
			ResponseBody: `{"slug": "academic-cant-edit-course", "details": "academic can't edit course"}`,
		},
		{
			Name:        "unexpected_error",
			RequestBody: `{"id": "8b41b8a4-4821-4029-87fc-a79ae0713cd5"}`,
			Authorized:  course.MustNewAcademic("60cddd22-8718-4d18-921b-8935f85ed7b8", course.TeacherType),
			CourseID:    "914c9a37-504c-496b-9715-f0ff2c8917ab",
			Command: app.AddCollaboratorCommand{
				CourseID:       "914c9a37-504c-496b-9715-f0ff2c8917ab",
				Academic:       course.MustNewAcademic("60cddd22-8718-4d18-921b-8935f85ed7b8", course.TeacherType),
				CollaboratorID: "8b41b8a4-4821-4029-87fc-a79ae0713cd5",
			},
			PrepareHandler: func(expectedCommand app.AddCollaboratorCommand) mock.AddCollaboratorsHandler {
				return func(_ context.Context, givenCommand app.AddCollaboratorCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return errors.New("unexpected error")
				}
			},
			StatusCode:   http.StatusInternalServerError,
			ResponseBody: `{"slug": "unexpected-error", "details": "unexpected error"}`,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			application := app.Application{
				Commands: app.Commands{
					AddCollaborator: c.PrepareHandler(c.Command),
				},
			}
			h := newHandler(t, application)

			w := httptest.NewRecorder()
			r := newRequest(
				t,
				http.MethodPut, fmt.Sprintf("/courses/%s/collaborators", c.CourseID),
				c.RequestBody, c.Authorized,
			)

			h.ServeHTTP(w, r)

			require.Equal(t, c.StatusCode, w.Code)

			if c.ResponseBody != "" {
				require.JSONEq(t, c.ResponseBody, w.Body.String())
			}
		})
	}
}
