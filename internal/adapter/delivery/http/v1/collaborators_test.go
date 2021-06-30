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

func TestHandler_RemoveCollaboratorFromCourse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		Authorized     course.Academic
		CourseID       string
		CollaboratorID string
		Command        app.RemoveCollaboratorCommand
		PrepareHandler func(expectedCommand app.RemoveCollaboratorCommand) mock.RemoveCollaboratorHandler
		StatusCode     int
		ResponseBody   string
	}{
		{
			Name:           "collaborator_removed_from_course",
			Authorized:     course.MustNewAcademic("bde6a399-1df0-45dd-8f75-a521503f2813", course.TeacherType),
			CourseID:       "58432a87-eda9-4c2c-bee3-ff5873a95708",
			CollaboratorID: "2c131b4a-e626-4dd2-a3ee-bbf2f20c025c",
			Command: app.RemoveCollaboratorCommand{
				CourseID:       "58432a87-eda9-4c2c-bee3-ff5873a95708",
				Academic:       course.MustNewAcademic("bde6a399-1df0-45dd-8f75-a521503f2813", course.TeacherType),
				CollaboratorID: "2c131b4a-e626-4dd2-a3ee-bbf2f20c025c",
			},
			PrepareHandler: func(expectedCommand app.RemoveCollaboratorCommand) mock.RemoveCollaboratorHandler {
				return func(_ context.Context, givenCommand app.RemoveCollaboratorCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return nil
				}
			},
			StatusCode:   http.StatusNoContent,
			ResponseBody: "",
		},
		{
			Name:           "course_not_found",
			Authorized:     course.MustNewAcademic("e957ab45-9f28-46b5-b2c1-1ec2943a6099", course.TeacherType),
			CourseID:       "2b906d15-b327-4ffc-ba91-b43cb20fe269",
			CollaboratorID: "b34ffb21-94d4-4af7-9faf-8df142ab8466",
			Command: app.RemoveCollaboratorCommand{
				CourseID:       "2b906d15-b327-4ffc-ba91-b43cb20fe269",
				Academic:       course.MustNewAcademic("e957ab45-9f28-46b5-b2c1-1ec2943a6099", course.TeacherType),
				CollaboratorID: "b34ffb21-94d4-4af7-9faf-8df142ab8466",
			},
			PrepareHandler: func(expectedCommand app.RemoveCollaboratorCommand) mock.RemoveCollaboratorHandler {
				return func(_ context.Context, givenCommand app.RemoveCollaboratorCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return app.ErrCourseDoesntExist
				}
			},
			StatusCode:   http.StatusNotFound,
			ResponseBody: `{"slug": "course-not-found", "details": "course doesn't exist"}`,
		},
		{
			Name:           "course_collaborator_not_found",
			Authorized:     course.MustNewAcademic("7bae715a-f7a0-467b-87af-5b3843453b10", course.TeacherType),
			CourseID:       "d3d3392f-547e-41ff-800d-2851b8701247",
			CollaboratorID: "e1f995ec-afed-4b27-a1d7-0fb74f04b363",
			Command: app.RemoveCollaboratorCommand{
				CourseID:       "d3d3392f-547e-41ff-800d-2851b8701247",
				Academic:       course.MustNewAcademic("7bae715a-f7a0-467b-87af-5b3843453b10", course.TeacherType),
				CollaboratorID: "e1f995ec-afed-4b27-a1d7-0fb74f04b363",
			},
			PrepareHandler: func(expectedCommand app.RemoveCollaboratorCommand) mock.RemoveCollaboratorHandler {
				return func(_ context.Context, givenCommand app.RemoveCollaboratorCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return course.ErrCourseHasNoSuchCollaborator
				}
			},
			StatusCode:   http.StatusNotFound,
			ResponseBody: `{"slug": "course-collaborator-not-found", "details": "course has no such collaborator"}`,
		},
		{
			Name:           "academic_cant_edit_course",
			Authorized:     course.MustNewAcademic("1dcb96be-9d83-43b0-8598-a0992d3bdfcd", course.TeacherType),
			CourseID:       "4a31f717-2bfb-46e2-8ec8-8c67ff6e9ad4",
			CollaboratorID: "a1d83105-7d26-47c2-9baa-b899652480fc",
			Command: app.RemoveCollaboratorCommand{
				CourseID:       "4a31f717-2bfb-46e2-8ec8-8c67ff6e9ad4",
				Academic:       course.MustNewAcademic("1dcb96be-9d83-43b0-8598-a0992d3bdfcd", course.TeacherType),
				CollaboratorID: "a1d83105-7d26-47c2-9baa-b899652480fc",
			},
			PrepareHandler: func(expectedCommand app.RemoveCollaboratorCommand) mock.RemoveCollaboratorHandler {
				return func(_ context.Context, givenCommand app.RemoveCollaboratorCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return course.AcademicCantEditCourseError{}
				}
			},
			StatusCode:   http.StatusForbidden,
			ResponseBody: `{"slug": "academic-cant-edit-course", "details": "academic can't edit course"}`,
		},
		{
			Name:           "unexpected_error",
			Authorized:     course.MustNewAcademic("fa2424e7-0ed9-49bb-9c29-fc8ee128da73", course.TeacherType),
			CourseID:       "4114b1fa-0854-439c-943a-6362926fbf15",
			CollaboratorID: "1549b8be-79aa-400a-ae8a-b7839491105b",
			Command: app.RemoveCollaboratorCommand{
				CourseID:       "4114b1fa-0854-439c-943a-6362926fbf15",
				Academic:       course.MustNewAcademic("fa2424e7-0ed9-49bb-9c29-fc8ee128da73", course.TeacherType),
				CollaboratorID: "1549b8be-79aa-400a-ae8a-b7839491105b",
			},
			PrepareHandler: func(expectedCommand app.RemoveCollaboratorCommand) mock.RemoveCollaboratorHandler {
				return func(_ context.Context, givenCommand app.RemoveCollaboratorCommand) error {
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
					RemoveCollaborator: c.PrepareHandler(c.Command),
				},
			}
			h := newHandler(t, application)

			w := httptest.NewRecorder()
			r := newRequest(
				t,
				http.MethodDelete, fmt.Sprintf("/courses/%s/collaborators/%s", c.CourseID, c.CollaboratorID),
				"", c.Authorized,
			)

			h.ServeHTTP(w, r)

			require.Equal(t, c.StatusCode, w.Code)

			if c.ResponseBody != "" {
				require.JSONEq(t, c.ResponseBody, w.Body.String())
			}
		})
	}
}
