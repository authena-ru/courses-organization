package v1_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/app/command/mock"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestHandler_AddTaskToCourse(t *testing.T) {
	t.Parallel()

	const courseID = "41542820-d331-4164-9384-1a51206cd8ce"

	testCases := []struct {
		Name               string
		RequestBody        string
		Authorized         course.Academic
		Command            app.AddTaskCommand
		PrepareHandler     func(expectedCommand app.AddTaskCommand) mock.AddTaskHandler
		StatusCode         int
		ResponseBody       string
		ExpectedTaskNumber int
	}{
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

			require.Equal(t, c.StatusCode, w.Code)
			if c.StatusCode == http.StatusCreated {
				require.Equal(
					t,
					fmt.Sprintf("/courses/%v/tasks/%v", courseID, c.ExpectedTaskNumber),
					w.Header().Get("Content-Location"),
				)
			}

			if c.ResponseBody != "" {
				require.JSONEq(t, c.RequestBody, w.Body.String())
			}
		})
	}
}
