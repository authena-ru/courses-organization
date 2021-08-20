package v1_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/app/command/mock"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestHandler_AddTaskToCourse(t *testing.T) {
	t.Parallel()

	const courseID = "41542820-d331-4164-9384-1a51206cd8ce"

	var (
		tooLongTitle      = strings.Repeat("x", 201)
		tooLongOutputData = strings.Repeat("∞", 1001)
	)

	testCases := []struct {
		Name                 string
		RequestBody          string
		Authorized           course.Academic
		Command              app.AddTaskCommand
		PrepareHandler       func(expectedCommand app.AddTaskCommand) mock.AddTaskHandler
		StatusCode           int
		ShouldBeResponseBody bool
		ResponseBody         string
		ExpectedTaskNumber   int
	}{
		{
			Name: "bad_request",
			RequestBody: `{
				"title": "Task title"",
				"description": "Task description",
				"type": "MANUAL_CHECKING",
			}`,
			Authorized: course.MustNewAcademic("0e248d6b-ce90-4588-8c24-51928a7de937", course.TeacherType),
			PrepareHandler: func(_ app.AddTaskCommand) mock.AddTaskHandler {
				return func(_ context.Context, _ app.AddTaskCommand) (int, error) {
					return 0, nil
				}
			},
			StatusCode:           http.StatusBadRequest,
			ShouldBeResponseBody: true,
			ResponseBody:         `{"slug":"bad-request", "details":"invalid character '\"' after object key:value pair"}`,
		},
		{
			Name: "manual_checking_task_added_to_course",
			RequestBody: `{
				"title": "Task title",
				"description": "Task description",
				"type": "MANUAL_CHECKING",
				"deadline": {
					"excellentGradeTime": "2021-09-03",
					"goodGradeTIme": "2021-09-24"
				}
			}`,
			Authorized: course.MustNewAcademic("458130df-25c7-4245-9a49-acd0e6c461b7", course.TeacherType),
			Command: app.AddTaskCommand{
				Academic:        course.MustNewAcademic("458130df-25c7-4245-9a49-acd0e6c461b7", course.TeacherType),
				CourseID:        courseID,
				TaskTitle:       "Task title",
				TaskDescription: "Task description",
				TaskType:        course.ManualCheckingType,
				Deadline: course.MustNewDeadline(
					time.Date(2021, time.September, 3, 0, 0, 0, 0, time.UTC),
					time.Date(2021, time.September, 24, 0, 0, 0, 0, time.UTC),
				),
			},
			PrepareHandler: func(expectedCommand app.AddTaskCommand) mock.AddTaskHandler {
				return func(_ context.Context, givenCommand app.AddTaskCommand) (int, error) {
					requireAddTaskCommandsEquals(t, expectedCommand, givenCommand)

					return 1, nil
				}
			},
			StatusCode:         http.StatusCreated,
			ExpectedTaskNumber: 1,
		},
		{
			Name: "testing_task_added_to_course",
			RequestBody: `{
				"title": "Testing task title",
				"description": "Testing task description",
				"type": "TESTING",
				"points": [
					{
						"description": "test point description", 
						"variants": ["Yes", "No"],
						"correctVariantNumbers": [0]
					}
				]
			}`,
			Authorized: course.MustNewAcademic("3f568bc5-8fc9-4535-ae06-d3cefcb0972c", course.TeacherType),
			Command: app.AddTaskCommand{
				Academic:        course.MustNewAcademic("3f568bc5-8fc9-4535-ae06-d3cefcb0972c", course.TeacherType),
				CourseID:        courseID,
				TaskTitle:       "Testing task title",
				TaskDescription: "Testing task description",
				TaskType:        course.TestingType,
				TestPoints: []course.TestPoint{
					course.MustNewTestPoint("test point description", []string{"Yes", "No"}, []int{0}),
				},
			},
			PrepareHandler: func(expectedCommand app.AddTaskCommand) mock.AddTaskHandler {
				return func(_ context.Context, givenCommand app.AddTaskCommand) (int, error) {
					requireAddTaskCommandsEquals(t, expectedCommand, givenCommand)

					return 2, nil
				}
			},
			StatusCode:         http.StatusCreated,
			ExpectedTaskNumber: 2,
		},
		{
			Name: "auto_code_checking_added_to_course",
			RequestBody: `{
				"title": "Auto code checking task title",
				"description": "Auto code checking task description",
				"type": "AUTO_CODE_CHECKING",
				"testData": [
					{
						"inputData": "1 + 1",
						"outputData": "2"
					}
				],
				"deadline": {
					"excellentGradeTime": "2021-09-01",
					"goodGradeTime": "2021-10-01"
				}
			}`,
			Authorized: course.MustNewAcademic("13a44c9c-c950-4a7e-a88b-4d6d9d5748f8", course.TeacherType),
			Command: app.AddTaskCommand{
				Academic:        course.MustNewAcademic("13a44c9c-c950-4a7e-a88b-4d6d9d5748f8", course.TeacherType),
				CourseID:        courseID,
				TaskTitle:       "Auto code checking task title",
				TaskDescription: "Auto code checking task description",
				TaskType:        course.AutoCodeCheckingType,
				TestData:        []course.TestData{course.MustNewTestData("1 + 1", "2")},
				Deadline: course.MustNewDeadline(
					time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC),
				),
			},
			PrepareHandler: func(expectedCommand app.AddTaskCommand) mock.AddTaskHandler {
				return func(_ context.Context, givenCommand app.AddTaskCommand) (int, error) {
					requireAddTaskCommandsEquals(t, expectedCommand, givenCommand)

					return 3, nil
				}
			},
			StatusCode:         http.StatusCreated,
			ExpectedTaskNumber: 3,
		},
		{
			Name: "course_not_found",
			RequestBody: `{
				"title": "Manual checking task #2",
				"description": "Manual checking task #2 description",
				"type": "MANUAL_CHECKING"
			}`,
			Authorized: course.MustNewAcademic("8140f5a7-7d3f-4cc2-8a68-d0bcb075ae32", course.TeacherType),
			Command: app.AddTaskCommand{
				Academic:        course.MustNewAcademic("8140f5a7-7d3f-4cc2-8a68-d0bcb075ae32", course.TeacherType),
				CourseID:        courseID,
				TaskTitle:       "Manual checking task #2",
				TaskDescription: "Manual checking task #2 description",
				TaskType:        course.ManualCheckingType,
			},
			PrepareHandler: func(expectedCommand app.AddTaskCommand) mock.AddTaskHandler {
				return func(_ context.Context, givenCommand app.AddTaskCommand) (int, error) {
					requireAddTaskCommandsEquals(t, expectedCommand, givenCommand)

					return 0, app.ErrCourseDoesntExist
				}
			},
			StatusCode:           http.StatusNotFound,
			ShouldBeResponseBody: true,
			ResponseBody:         `{"slug": "course-not-found", "details": "course doesn't exist"}`,
		},
		{
			Name: "invalid_task_parameters",
			RequestBody: fmt.Sprintf(`{
				"title": "%s",
				"description": "Description for task with too long title",
				"type": "MANUAL_CHECKING"
			}`, tooLongTitle),
			Authorized: course.MustNewAcademic("e9ab09c7-fa8c-47b4-af0b-9f8934f81488", course.TeacherType),
			Command: app.AddTaskCommand{
				Academic:        course.MustNewAcademic("e9ab09c7-fa8c-47b4-af0b-9f8934f81488", course.TeacherType),
				CourseID:        courseID,
				TaskTitle:       tooLongTitle,
				TaskDescription: "Description for task with too long title",
				TaskType:        course.ManualCheckingType,
			},
			PrepareHandler: func(expectedCommand app.AddTaskCommand) mock.AddTaskHandler {
				return func(_ context.Context, givenCommand app.AddTaskCommand) (int, error) {
					requireAddTaskCommandsEquals(t, expectedCommand, givenCommand)

					return 0, course.ErrTaskTitleTooLong
				}
			},
			StatusCode:           http.StatusUnprocessableEntity,
			ShouldBeResponseBody: true,
			ResponseBody:         `{"slug": "invalid-task-parameters", "details": "task title too long"}`,
		},
		{
			Name: "invalid_deadline",
			RequestBody: `{
				"title": "Some task #69",
				"description": "69",
				"type": "MANUAL_CHECKING",
				"deadline": {
					"excellentGradeTime": "2021-03-01",
					"goodGradeTime": "2021-02-01"
				}
			}`,
			Authorized: course.MustNewAcademic("f695f54c-65e0-46da-a3ce-ffe93a13641b", course.TeacherType),
			PrepareHandler: func(_ app.AddTaskCommand) mock.AddTaskHandler {
				return func(_ context.Context, givenCommand app.AddTaskCommand) (int, error) {
					return 0, nil
				}
			},
			StatusCode:           http.StatusUnprocessableEntity,
			ShouldBeResponseBody: true,
			ResponseBody:         `{"slug": "invalid-deadline", "details": "excellent grade time after good"}`,
		},
		{
			Name: "invalid_test_data",
			RequestBody: fmt.Sprintf(`{
				"title": "Test task",
				"description": "Test task description",
				"type": "AUTO_CODE_CHECKING",
				"testData": [
					{
						"inputData": "2 * 2",
						"outputData": "4"
					},
					{
						"inputData": "inf + inf",
						"outputData": "%s"
					}
				]
			}`, tooLongOutputData),
			Authorized: course.MustNewAcademic("c29e0dc8-59f4-48b0-926e-281bd7ee56b8", course.TeacherType),
			PrepareHandler: func(_ app.AddTaskCommand) mock.AddTaskHandler {
				return func(_ context.Context, _ app.AddTaskCommand) (int, error) {
					return 0, nil
				}
			},
			StatusCode:           http.StatusUnprocessableEntity,
			ShouldBeResponseBody: true,
			ResponseBody:         `{"slug": "invalid-test-data", "details": "test output data too long"}`,
		},
		{
			Name: "invalid_test_point",
			RequestBody: `{
				"title": "TESTING task",
				"description": "TESTING task description",
				"type": "TESTING",
				"points": [
					{
						"description": "? ? ?",
						"variants": ["?", "? ?", "? ? ?"],
						"correctVariantNumbers": [3]
					}
				]
			}`,
			Authorized: course.MustNewAcademic("fc0601f7-e8b2-4a0b-8adc-38d82eb4f80d", course.TeacherType),
			PrepareHandler: func(_ app.AddTaskCommand) mock.AddTaskHandler {
				return func(_ context.Context, _ app.AddTaskCommand) (int, error) {
					return 0, nil
				}
			},
			StatusCode:           http.StatusUnprocessableEntity,
			ShouldBeResponseBody: true,
			ResponseBody:         `{"slug": "invalid-test-point", "details": "invalid test point variant number"}`,
		},
		{
			Name: "invalid_task_type",
			RequestBody: `{
				"title": "Task without type",
				"description": "...Description...",
				"type": "UNKNOWN"
			}`,
			Authorized: course.MustNewAcademic("920e5b80-b7d2-468f-8fdd-707650ff16f2", course.TeacherType),
			PrepareHandler: func(expectedCommand app.AddTaskCommand) mock.AddTaskHandler {
				return func(_ context.Context, _ app.AddTaskCommand) (int, error) {
					return 0, nil
				}
			},
			StatusCode:           http.StatusUnprocessableEntity,
			ShouldBeResponseBody: true,
			ResponseBody:         `{"slug": "invalid-task-type", "details": ""}`,
		},
		{
			Name: "academic_cant_edit_course",
			RequestBody: `{
				"title": "Some task",
				"description": "some description",
				"type": "MANUAL_CHECKING"
			}`,
			Authorized: course.MustNewAcademic("b4510051-be17-4fd6-857b-088d6de3cbab", course.StudentType),
			Command: app.AddTaskCommand{
				Academic:        course.MustNewAcademic("b4510051-be17-4fd6-857b-088d6de3cbab", course.StudentType),
				CourseID:        courseID,
				TaskTitle:       "Some task",
				TaskDescription: "some description",
				TaskType:        course.ManualCheckingType,
			},
			PrepareHandler: func(expectedCommand app.AddTaskCommand) mock.AddTaskHandler {
				return func(_ context.Context, givenCommand app.AddTaskCommand) (int, error) {
					requireAddTaskCommandsEquals(t, expectedCommand, givenCommand)

					return 0, course.AcademicCantEditCourseError{}
				}
			},
			StatusCode:           http.StatusForbidden,
			ShouldBeResponseBody: true,
			ResponseBody:         `{"slug": "academic-cant-edit-course", "details": "academic can't edit course"}`,
		},
		{
			Name: "unexpected-error",
			RequestBody: `{
				"title": "T A S K",
				"description": "D E S C R I P T I O N",
				"type": "MANUAL_CHECKING"
			}`,
			Authorized: course.MustNewAcademic("ea014b8f-894b-4b6b-a810-3a9872581ad2", course.TeacherType),
			Command: app.AddTaskCommand{
				Academic:        course.MustNewAcademic("ea014b8f-894b-4b6b-a810-3a9872581ad2", course.TeacherType),
				CourseID:        courseID,
				TaskTitle:       "T A S K",
				TaskDescription: "D E S C R I P T I O N",
				TaskType:        course.ManualCheckingType,
			},
			PrepareHandler: func(expectedCommand app.AddTaskCommand) mock.AddTaskHandler {
				return func(ctx context.Context, givenCommand app.AddTaskCommand) (int, error) {
					requireAddTaskCommandsEquals(t, expectedCommand, givenCommand)

					return 0, errors.New("unexpected error")
				}
			},
			StatusCode:           http.StatusInternalServerError,
			ShouldBeResponseBody: true,
			ResponseBody:         `{"slug": "unexpected-error", "details": "unexpected error"}`,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			application := app.Application{
				Commands: app.Commands{
					AddTask: c.PrepareHandler(c.Command),
				},
			}
			h := newHTTPHandler(t, application)

			w := httptest.NewRecorder()
			r := newHTTPRequest(
				t,
				http.MethodPost, fmt.Sprintf("/courses/%v/tasks", courseID),
				c.RequestBody, c.Authorized,
			)

			h.ServeHTTP(w, r)

			require.Equalf(t, c.StatusCode, w.Code, "status codes are not equal")
			if c.StatusCode == http.StatusCreated {
				require.Equal(
					t,
					fmt.Sprintf("/courses/%v/tasks/%v", courseID, c.ExpectedTaskNumber),
					w.Header().Get("Content-Location"),
				)
			}

			if c.ShouldBeResponseBody {
				require.JSONEq(t, c.ResponseBody, w.Body.String())
			}
		})
	}
}
