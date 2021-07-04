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

func TestHandler_AddStudentToCourse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		RequestBody    string
		Authorized     course.Academic
		CourseID       string
		Command        app.AddStudentCommand
		PrepareHandler func(expectedCommand app.AddStudentCommand) mock.AddStudentHandler
		StatusCode     int
		ResponseBody   string
	}{
		{
			Name:        "student_added_to_course",
			RequestBody: `{"id": "473d6e1c-2d3d-4a40-9a95-4037edaa33e3"}`,
			Authorized:  course.MustNewAcademic("70e8ed77-09d2-43c1-82da-deacd58facd6", course.TeacherType),
			CourseID:    "2e5d112c-987c-4ac1-8c54-9ad1e54eb408",
			Command: app.AddStudentCommand{
				Academic:  course.MustNewAcademic("70e8ed77-09d2-43c1-82da-deacd58facd6", course.TeacherType),
				CourseID:  "2e5d112c-987c-4ac1-8c54-9ad1e54eb408",
				StudentID: "473d6e1c-2d3d-4a40-9a95-4037edaa33e3",
			},
			PrepareHandler: func(expectedCommand app.AddStudentCommand) mock.AddStudentHandler {
				return func(_ context.Context, givenCommand app.AddStudentCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return nil
				}
			},
			StatusCode:   http.StatusNoContent,
			ResponseBody: "",
		},
		{
			Name:        "bad_request",
			RequestBody: `{"id": 123434}`,
			Authorized:  course.MustNewAcademic("78a55474-c7c6-4677-acc5-f9b76376c4ab", course.TeacherType),
			CourseID:    "922b99f5-e1a8-48ef-a915-cc9f48bd4512",
			PrepareHandler: func(expectedCommand app.AddStudentCommand) mock.AddStudentHandler {
				return func(_ context.Context, _ app.AddStudentCommand) error {
					return nil
				}
			},
			StatusCode: http.StatusBadRequest,
		},
		{
			Name:        "course_not_found",
			RequestBody: `{"id": "7a91cda0-6097-44ae-a9c8-286a0121c6c2"}`,
			Authorized:  course.MustNewAcademic("8994a75f-9706-4a95-8404-4fb894c5b23a", course.TeacherType),
			CourseID:    "3e0126e1-891e-4813-b65f-4ab8147efdc9",
			Command: app.AddStudentCommand{
				Academic:  course.MustNewAcademic("8994a75f-9706-4a95-8404-4fb894c5b23a", course.TeacherType),
				CourseID:  "3e0126e1-891e-4813-b65f-4ab8147efdc9",
				StudentID: "7a91cda0-6097-44ae-a9c8-286a0121c6c2",
			},
			PrepareHandler: func(expectedCommand app.AddStudentCommand) mock.AddStudentHandler {
				return func(_ context.Context, givenCommand app.AddStudentCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return app.ErrCourseDoesntExist
				}
			},
			StatusCode:   http.StatusNotFound,
			ResponseBody: `{"slug": "course-not-found", "details": "course doesn't exist"}`,
		},
		{
			Name:        "student_not_found",
			RequestBody: `{"id": "e3b527cf-7acb-4b92-8773-b706eea40efa"}`,
			Authorized:  course.MustNewAcademic("f83dc0a4-8839-44f5-9b5f-e21e29ff7b90", course.TeacherType),
			CourseID:    "20757653-0def-4329-aac5-f80d5aff9f99",
			Command: app.AddStudentCommand{
				Academic:  course.MustNewAcademic("f83dc0a4-8839-44f5-9b5f-e21e29ff7b90", course.TeacherType),
				CourseID:  "20757653-0def-4329-aac5-f80d5aff9f99",
				StudentID: "e3b527cf-7acb-4b92-8773-b706eea40efa",
			},
			PrepareHandler: func(expectedCommand app.AddStudentCommand) mock.AddStudentHandler {
				return func(_ context.Context, givenCommand app.AddStudentCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return app.ErrStudentDoesntExist
				}
			},
			StatusCode:   http.StatusUnprocessableEntity,
			ResponseBody: `{"slug": "student-not-found", "details": "student doesn't exist"}`,
		},
		{
			Name:        "academic_cant_edit_course",
			RequestBody: `{"id": "6f690f76-6e86-4201-a56b-d0ee270a9928"}`,
			Authorized:  course.MustNewAcademic("ab3f3c6e-a2bf-4776-a0df-91941e62f1c7", course.TeacherType),
			CourseID:    "e3f1b7cc-6c64-4454-9534-ffe5d254a7de",
			Command: app.AddStudentCommand{
				Academic:  course.MustNewAcademic("ab3f3c6e-a2bf-4776-a0df-91941e62f1c7", course.TeacherType),
				CourseID:  "e3f1b7cc-6c64-4454-9534-ffe5d254a7de",
				StudentID: "6f690f76-6e86-4201-a56b-d0ee270a9928",
			},
			PrepareHandler: func(expectedCommand app.AddStudentCommand) mock.AddStudentHandler {
				return func(_ context.Context, givenCommand app.AddStudentCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return course.AcademicCantEditCourseError{}
				}
			},
			StatusCode:   http.StatusForbidden,
			ResponseBody: `{"slug": "academic-cant-edit-course", "details": "academic can't edit course"}`,
		},
		{
			Name:        "unexpected_error",
			RequestBody: `{"id": "9ca7f110-ff81-4f11-a795-fa853c5aabcc"}`,
			Authorized:  course.MustNewAcademic("ec24ea6f-3393-4372-80c7-77609450de9b", course.TeacherType),
			CourseID:    "a913ad8a-e892-435d-8d91-5ca78fc8488d",
			Command: app.AddStudentCommand{
				Academic:  course.MustNewAcademic("ec24ea6f-3393-4372-80c7-77609450de9b", course.TeacherType),
				CourseID:  "a913ad8a-e892-435d-8d91-5ca78fc8488d",
				StudentID: "9ca7f110-ff81-4f11-a795-fa853c5aabcc",
			},
			PrepareHandler: func(expectedCommand app.AddStudentCommand) mock.AddStudentHandler {
				return func(_ context.Context, givenCommand app.AddStudentCommand) error {
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
					AddStudent: c.PrepareHandler(c.Command),
				},
			}
			h := newHTTPHandler(t, application)

			w := httptest.NewRecorder()
			r := newHTTPRequest(
				t,
				http.MethodPut, fmt.Sprintf("/courses/%s/students", c.CourseID),
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

func TestHandler_RemoveStudentFromCourse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		Authorized     course.Academic
		CourseID       string
		StudentID      string
		Command        app.RemoveStudentCommand
		PrepareHandler func(expectedCommand app.RemoveStudentCommand) mock.RemoveStudentHandler
		StatusCode     int
		ResponseBody   string
	}{
		{
			Name:       "student_removed_from_course",
			Authorized: course.MustNewAcademic("4eff6499-fbde-4090-a552-293be8053e53", course.TeacherType),
			CourseID:   "8c4d464c-420e-4eb3-9bfa-94cf8a97883a",
			StudentID:  "7ef909eb-38be-4531-9055-2ed3ad79a260",
			Command: app.RemoveStudentCommand{
				Academic:  course.MustNewAcademic("4eff6499-fbde-4090-a552-293be8053e53", course.TeacherType),
				CourseID:  "8c4d464c-420e-4eb3-9bfa-94cf8a97883a",
				StudentID: "7ef909eb-38be-4531-9055-2ed3ad79a260",
			},
			PrepareHandler: func(expectedCommand app.RemoveStudentCommand) mock.RemoveStudentHandler {
				return func(_ context.Context, givenCommand app.RemoveStudentCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return nil
				}
			},
			StatusCode:   http.StatusNoContent,
			ResponseBody: "",
		},
		{
			Name:       "course_not_found",
			Authorized: course.MustNewAcademic("76d00409-08b1-4c3e-b81b-2a0bbb73683e", course.TeacherType),
			CourseID:   "8052b540-4189-4bff-9035-10695891623a",
			StudentID:  "7132d45f-ca7f-4bc4-82e1-169f9443738e",
			Command: app.RemoveStudentCommand{
				Academic:  course.MustNewAcademic("76d00409-08b1-4c3e-b81b-2a0bbb73683e", course.TeacherType),
				CourseID:  "8052b540-4189-4bff-9035-10695891623a",
				StudentID: "7132d45f-ca7f-4bc4-82e1-169f9443738e",
			},
			PrepareHandler: func(expectedCommand app.RemoveStudentCommand) mock.RemoveStudentHandler {
				return func(_ context.Context, givenCommand app.RemoveStudentCommand) error {
					return app.ErrCourseDoesntExist
				}
			},
			StatusCode:   http.StatusNotFound,
			ResponseBody: `{"slug": "course-not-found", "details": "course doesn't exist"}`,
		},
		{
			Name:       "course_student_not_found",
			Authorized: course.MustNewAcademic("f888e716-35f1-479b-9565-d3720c289508", course.TeacherType),
			CourseID:   "ae0e3b76-39dc-46ed-9f01-72ea469ae33b",
			StudentID:  "91bc4990-5fd8-4b0b-b930-e2fea7266d9c",
			Command: app.RemoveStudentCommand{
				Academic:  course.MustNewAcademic("f888e716-35f1-479b-9565-d3720c289508", course.TeacherType),
				CourseID:  "ae0e3b76-39dc-46ed-9f01-72ea469ae33b",
				StudentID: "91bc4990-5fd8-4b0b-b930-e2fea7266d9c",
			},
			PrepareHandler: func(expectedCommand app.RemoveStudentCommand) mock.RemoveStudentHandler {
				return func(_ context.Context, givenCommand app.RemoveStudentCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return course.ErrCourseHasNoSuchStudent
				}
			},
			StatusCode:   http.StatusNotFound,
			ResponseBody: `{"slug": "course-student-not-found", "details": "course has no such student"}`,
		},
		{
			Name:       "academic_cant_edit_course",
			Authorized: course.MustNewAcademic("b3917d2e-e3fe-4915-a023-96aa499d9bd1", course.StudentType),
			CourseID:   "e8b5975d-d87e-42fe-b305-068bf62b85b2",
			StudentID:  "fac6cee9-f665-4b11-bb8e-33a9185e0b39",
			Command: app.RemoveStudentCommand{
				Academic:  course.MustNewAcademic("b3917d2e-e3fe-4915-a023-96aa499d9bd1", course.StudentType),
				CourseID:  "e8b5975d-d87e-42fe-b305-068bf62b85b2",
				StudentID: "fac6cee9-f665-4b11-bb8e-33a9185e0b39",
			},
			PrepareHandler: func(expectedCommand app.RemoveStudentCommand) mock.RemoveStudentHandler {
				return func(_ context.Context, givenCommand app.RemoveStudentCommand) error {
					require.Equal(t, expectedCommand, givenCommand)

					return course.AcademicCantEditCourseError{}
				}
			},
			StatusCode:   http.StatusForbidden,
			ResponseBody: `{"slug": "academic-cant-edit-course", "details": "academic can't edit course"}`,
		},
		{
			Name:       "unexpected_error",
			Authorized: course.MustNewAcademic("d9d01cb4-a41c-47cb-be40-0bf5d350143d", course.TeacherType),
			CourseID:   "413e271b-6d02-436a-9cde-e8e52ebeb10c",
			StudentID:  "0005bf02-2851-4cb0-8fcc-4bbf81b00b94",
			Command: app.RemoveStudentCommand{
				Academic:  course.MustNewAcademic("d9d01cb4-a41c-47cb-be40-0bf5d350143d", course.TeacherType),
				CourseID:  "413e271b-6d02-436a-9cde-e8e52ebeb10c",
				StudentID: "0005bf02-2851-4cb0-8fcc-4bbf81b00b94",
			},
			PrepareHandler: func(expectedCommand app.RemoveStudentCommand) mock.RemoveStudentHandler {
				return func(_ context.Context, givenCommand app.RemoveStudentCommand) error {
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
					RemoveStudent: c.PrepareHandler(c.Command),
				},
			}
			h := newHTTPHandler(t, application)

			w := httptest.NewRecorder()
			r := newHTTPRequest(
				t,
				http.MethodDelete, fmt.Sprintf("/courses/%s/students/%s",
					c.CourseID, c.StudentID), "", c.Authorized,
			)

			h.ServeHTTP(w, r)

			require.Equal(t, c.StatusCode, w.Code)

			if c.ResponseBody != "" {
				require.JSONEq(t, c.ResponseBody, w.Body.String())
			}
		})
	}
}
