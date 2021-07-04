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

func TestHandler_CreateCourse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name            string
		RequestBody     string
		Authorized      course.Academic
		Command         app.CreateCourseCommand
		PrepareHandler  func(expectedCommand app.CreateCourseCommand) mock.CreateCourseHandler
		StatusCode      int
		ResponseBody    string
		ContentLocation string
	}{
		{
			Name:        "course_created",
			RequestBody: `{"title": "New brand course", "started": true, "period": {"academicStartYear": 2020, "academicEndYear": 2021, "semester": "FIRST"}}`,
			Authorized:  course.MustNewAcademic("abb23a78-e3f7-4636-a14c-9bb5ba0fad42", course.TeacherType),
			Command: app.CreateCourseCommand{
				Academic:      course.MustNewAcademic("abb23a78-e3f7-4636-a14c-9bb5ba0fad42", course.TeacherType),
				CourseTitle:   "New brand course",
				CourseStarted: true,
				CoursePeriod:  course.MustNewPeriod(2020, 2021, course.FirstSemester),
			},
			PrepareHandler: func(expectedCommand app.CreateCourseCommand) mock.CreateCourseHandler {
				return func(_ context.Context, givenCommand app.CreateCourseCommand) (string, error) {
					require.Equal(t, expectedCommand, givenCommand)

					return "47e3a14f-c48d-4d86-b587-b4b319f4f733", nil
				}
			},
			StatusCode:      http.StatusCreated,
			ResponseBody:    "",
			ContentLocation: "/courses/47e3a14f-c48d-4d86-b587-b4b319f4f733",
		},
		{
			Name:        "bad_request",
			RequestBody: `{"title": 1}`,
			Authorized:  course.MustNewAcademic("486cc779-d99d-4f6c-bbd8-a59ceb496c23", course.TeacherType),
			PrepareHandler: func(expectedCommand app.CreateCourseCommand) mock.CreateCourseHandler {
				return func(_ context.Context, _ app.CreateCourseCommand) (string, error) {
					return "", nil
				}
			},
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"slug": "bad-request", "details": "json: cannot unmarshal number into Go struct field CreateCourseRequest.title of type string"}`,
		},
		{
			Name:        "invalid_course_parameters",
			RequestBody: `{"title": "", "started": false, "period": {"academicStartYear": 2021, "academicEndYear": 2022, "semester": "SECOND"}}`,
			Authorized:  course.MustNewAcademic("8e2f7b8e-7ad7-46e0-b62f-4a01c133ae8e", course.TeacherType),
			Command: app.CreateCourseCommand{
				Academic:      course.MustNewAcademic("8e2f7b8e-7ad7-46e0-b62f-4a01c133ae8e", course.TeacherType),
				CourseTitle:   "",
				CourseStarted: false,
				CoursePeriod:  course.MustNewPeriod(2021, 2022, course.SecondSemester),
			},
			PrepareHandler: func(expectedCommand app.CreateCourseCommand) mock.CreateCourseHandler {
				return func(_ context.Context, givenCommand app.CreateCourseCommand) (string, error) {
					require.Equal(t, expectedCommand, givenCommand)

					return "", course.ErrEmptyCourseTitle
				}
			},
			StatusCode:   http.StatusUnprocessableEntity,
			ResponseBody: `{"slug": "invalid-course-parameters", "details": "empty course title"}`,
		},
		{
			Name:        "not_teacher_cant_create_course",
			RequestBody: `{"title": "Mocked", "period": {"academicStartYear": 2022, "academicEndYear": 2023, "semester": "FIRST"}}`,
			Authorized:  course.MustNewAcademic("08136131-7ba6-4af1-a86d-7e10e43507eb", course.StudentType),
			Command: app.CreateCourseCommand{
				Academic:      course.MustNewAcademic("08136131-7ba6-4af1-a86d-7e10e43507eb", course.StudentType),
				CourseTitle:   "Mocked",
				CourseStarted: false,
				CoursePeriod:  course.MustNewPeriod(2022, 2023, course.FirstSemester),
			},
			PrepareHandler: func(expectedCommand app.CreateCourseCommand) mock.CreateCourseHandler {
				return func(_ context.Context, givenCommand app.CreateCourseCommand) (string, error) {
					require.Equal(t, expectedCommand, givenCommand)

					return "", course.ErrNotTeacherCantCreateCourse
				}
			},
			StatusCode:   http.StatusForbidden,
			ResponseBody: `{"slug": "not-teacher-cant-create-course", "details": "not teacher can't create course"}`,
		},
		{
			Name:        "unexpected_error",
			RequestBody: `{"title": "New mocked", "period": {"academicStartYear": 2024, "academicEndYear": 2025, "semester": "FIRST"}}`,
			Authorized:  course.MustNewAcademic("28897c1a-04eb-4b4c-a371-917578b188d0", course.TeacherType),
			Command: app.CreateCourseCommand{
				Academic:     course.MustNewAcademic("28897c1a-04eb-4b4c-a371-917578b188d0", course.TeacherType),
				CourseTitle:  "New mocked",
				CoursePeriod: course.MustNewPeriod(2024, 2025, course.FirstSemester),
			},
			PrepareHandler: func(expectedCommand app.CreateCourseCommand) mock.CreateCourseHandler {
				return func(_ context.Context, givenCommand app.CreateCourseCommand) (string, error) {
					require.Equal(t, expectedCommand, givenCommand)

					return "", errors.New("unexpected error")
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
					CreateCourse: c.PrepareHandler(c.Command),
				},
			}
			h := newHTTPHandler(t, application)

			w := httptest.NewRecorder()
			r := newHTTPRequest(t, http.MethodPost, "/courses", c.RequestBody, c.Authorized)

			h.ServeHTTP(w, r)

			require.Equal(t, c.StatusCode, w.Code)
			require.Equal(t, c.ContentLocation, w.Header().Get("Content-Location"))

			if c.ResponseBody != "" {
				require.JSONEq(t, c.ResponseBody, w.Body.String())
			}
		})
	}
}

func TestHandler_ExtendCourse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name            string
		RequestBody     string
		OriginCourseID  string
		Authorized      course.Academic
		Command         app.ExtendCourseCommand
		PrepareHandler  func(expectedCommand app.ExtendCourseCommand) mock.ExtendCourseHandler
		StatusCode      int
		ResponseBody    string
		ContentLocation string
	}{
		{
			Name:           "course_extended",
			RequestBody:    `{"title": "Extended", "started": true}`,
			OriginCourseID: "27a716e4-0a63-44d5-a549-73c5fc3be4df",
			Authorized:     course.MustNewAcademic("71d6d593-71e9-4507-bd12-cbe76bd1d339", course.TeacherType),
			Command: app.ExtendCourseCommand{
				Academic:       course.MustNewAcademic("71d6d593-71e9-4507-bd12-cbe76bd1d339", course.TeacherType),
				OriginCourseID: "27a716e4-0a63-44d5-a549-73c5fc3be4df",
				CourseTitle:    "Extended",
				CourseStarted:  true,
			},
			PrepareHandler: func(expectedCommand app.ExtendCourseCommand) mock.ExtendCourseHandler {
				return func(_ context.Context, givenCommand app.ExtendCourseCommand) (string, error) {
					require.Equal(t, expectedCommand, givenCommand)

					return "83d7502c-6327-40cc-a1ca-f6f44889137f", nil
				}
			},
			StatusCode:      http.StatusCreated,
			ResponseBody:    "",
			ContentLocation: "/courses/83d7502c-6327-40cc-a1ca-f6f44889137f",
		},
		{
			Name:           "bad_request",
			RequestBody:    `{"title": "", "started": ""}`,
			OriginCourseID: "3d851349-1d6d-43ce-aed9-ecac971512aa",
			Authorized:     course.MustNewAcademic("a45dc626-d478-45f2-b1ec-0268c61acbf7", course.TeacherType),
			PrepareHandler: func(expectedCommand app.ExtendCourseCommand) mock.ExtendCourseHandler {
				return func(_ context.Context, _ app.ExtendCourseCommand) (string, error) {
					return "", nil
				}
			},
			StatusCode:   http.StatusBadRequest,
			ResponseBody: `{"slug": "bad-request", "details": "json: cannot unmarshal string into Go struct field ExtendCourseRequest.started of type bool"}`,
		},
		{
			Name:           "course_not_found",
			RequestBody:    `{"title": "Extended mocked course"}`,
			OriginCourseID: "a2db2f55-20f3-43c2-b9c8-039d0ac5ce4e",
			Authorized:     course.MustNewAcademic("0b57e5bb-c66e-4da7-b6ea-ca6cd03ec46a", course.TeacherType),
			Command: app.ExtendCourseCommand{
				Academic:       course.MustNewAcademic("0b57e5bb-c66e-4da7-b6ea-ca6cd03ec46a", course.TeacherType),
				OriginCourseID: "a2db2f55-20f3-43c2-b9c8-039d0ac5ce4e",
				CourseTitle:    "Extended mocked course",
			},
			PrepareHandler: func(expectedCommand app.ExtendCourseCommand) mock.ExtendCourseHandler {
				return func(_ context.Context, givenCommand app.ExtendCourseCommand) (string, error) {
					require.Equal(t, expectedCommand, givenCommand)

					return "", app.ErrCourseDoesntExist
				}
			},
			StatusCode:   http.StatusNotFound,
			ResponseBody: `{"slug": "course-not-found", "details": "course doesn't exist"}`,
		},
		{
			Name:           "invalid_course_parameters",
			RequestBody:    `{"title": ""}`,
			OriginCourseID: "ecff625e-3101-4032-8372-4b04cc298548",
			Authorized:     course.MustNewAcademic("55c5cd7d-a5bc-4266-b6d6-ae81a5c30a5d", course.TeacherType),
			Command: app.ExtendCourseCommand{
				Academic:       course.MustNewAcademic("55c5cd7d-a5bc-4266-b6d6-ae81a5c30a5d", course.TeacherType),
				OriginCourseID: "ecff625e-3101-4032-8372-4b04cc298548",
				CourseTitle:    "",
			},
			PrepareHandler: func(expectedCommand app.ExtendCourseCommand) mock.ExtendCourseHandler {
				return func(_ context.Context, givenCommand app.ExtendCourseCommand) (string, error) {
					require.Equal(t, expectedCommand, givenCommand)

					return "", course.ErrEmptyCourseTitle
				}
			},
			StatusCode:   http.StatusUnprocessableEntity,
			ResponseBody: `{"slug": "invalid-course-parameters", "details": "empty course title"}`,
		},
		{
			Name:           "academic_cant_edit_course",
			RequestBody:    `{"title": "New"}`,
			OriginCourseID: "6c3c870c-b65f-448c-9f57-b735f5c595d9",
			Authorized:     course.MustNewAcademic("dfc3839c-28e4-4470-a183-f49699962ec7", course.TeacherType),
			Command: app.ExtendCourseCommand{
				Academic:       course.MustNewAcademic("dfc3839c-28e4-4470-a183-f49699962ec7", course.TeacherType),
				OriginCourseID: "6c3c870c-b65f-448c-9f57-b735f5c595d9",
				CourseTitle:    "New",
			},
			PrepareHandler: func(expectedCommand app.ExtendCourseCommand) mock.ExtendCourseHandler {
				return func(_ context.Context, givenCommand app.ExtendCourseCommand) (string, error) {
					require.Equal(t, expectedCommand, givenCommand)

					return "", course.AcademicCantEditCourseError{}
				}
			},
			StatusCode:   http.StatusForbidden,
			ResponseBody: `{"slug": "academic-cant-edit-course", "details": "academic can't edit course"}`,
		},
		{
			Name:           "not_teacher_cant_create_course",
			RequestBody:    `{"title": "New extended"}`,
			OriginCourseID: "398f1034-ba06-49c3-8c55-4f7ce81838c0",
			Authorized:     course.MustNewAcademic("1de58d6d-3179-4da4-9012-e4d8dee5f12c", course.StudentType),
			Command: app.ExtendCourseCommand{
				Academic:       course.MustNewAcademic("1de58d6d-3179-4da4-9012-e4d8dee5f12c", course.StudentType),
				OriginCourseID: "398f1034-ba06-49c3-8c55-4f7ce81838c0",
				CourseTitle:    "New extended",
			},
			PrepareHandler: func(expectedCommand app.ExtendCourseCommand) mock.ExtendCourseHandler {
				return func(_ context.Context, givenCommand app.ExtendCourseCommand) (string, error) {
					require.Equal(t, expectedCommand, givenCommand)

					return "", course.ErrNotTeacherCantCreateCourse
				}
			},
			StatusCode:   http.StatusForbidden,
			ResponseBody: `{"slug": "not-teacher-cant-create-course", "details": "not teacher can't create course"}`,
		},
		{
			Name:           "unexpected_error",
			RequestBody:    `{"title": "New brand extended"}`,
			OriginCourseID: "73786011-74c1-4e5c-abff-c667fab61454",
			Authorized:     course.MustNewAcademic("d0ff62b2-be40-43d5-b54b-103039a17b34", course.TeacherType),
			Command: app.ExtendCourseCommand{
				Academic:       course.MustNewAcademic("d0ff62b2-be40-43d5-b54b-103039a17b34", course.TeacherType),
				OriginCourseID: "73786011-74c1-4e5c-abff-c667fab61454",
				CourseTitle:    "New brand extended",
			},
			PrepareHandler: func(expectedCommand app.ExtendCourseCommand) mock.ExtendCourseHandler {
				return func(_ context.Context, givenCommand app.ExtendCourseCommand) (string, error) {
					require.Equal(t, expectedCommand, givenCommand)

					return "", errors.New("unexpected error")
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
					ExtendCourse: c.PrepareHandler(c.Command),
				},
			}
			h := newHTTPHandler(t, application)

			w := httptest.NewRecorder()
			r := newHTTPRequest(
				t,
				http.MethodPost, fmt.Sprintf("/courses/%s/extended", c.OriginCourseID),
				c.RequestBody, c.Authorized,
			)

			h.ServeHTTP(w, r)

			require.Equal(t, c.StatusCode, w.Code)
			require.Equal(t, c.ContentLocation, w.Header().Get("Content-Location"))

			if c.ResponseBody != "" {
				require.JSONEq(t, c.ResponseBody, w.Body.String())
			}
		})
	}
}
