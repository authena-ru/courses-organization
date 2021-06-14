package command_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/app/command"
	"github.com/authena-ru/courses-organization/internal/app/command/mock"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestAddTaskHandler_Handle(t *testing.T) {
	t.Parallel()
	var (
		courseID = "course-id"
		creator  = course.MustNewAcademic("creator-id", course.TeacherType)
	)
	addCourse := func(crs *course.Course) *mock.CoursesRepository {
		return mock.NewCoursesRepository(crs)
	}
	testCases := []struct {
		Name                     string
		Command                  command.AddTaskCommand
		PrepareCoursesRepository func(crs *course.Course) *mock.CoursesRepository
		IsErr                    func(err error) bool
	}{
		{
			Name: "add_manual_checking_task",
			Command: command.AddTaskCommand{
				Academic:        creator,
				CourseID:        courseID,
				TaskTitle:       "Manual checking task",
				TaskDescription: "Do this task",
				TaskType:        course.ManualChecking,
				Deadline: course.MustNewDeadline(
					time.Date(2043, time.November, 10, 0, 0, 0, 0, time.Local),
					time.Date(2043, time.November, 23, 0, 0, 0, 0, time.Local),
				),
			},
			PrepareCoursesRepository: addCourse,
		},
		{
			Name: "add_auto_code_checking_task",
			Command: command.AddTaskCommand{
				Academic:        creator,
				CourseID:        courseID,
				TaskTitle:       "Auto code checking task",
				TaskDescription: "Do this task",
				TaskType:        course.AutoCodeChecking,
				TestData:        []course.TestData{course.MustNewTestData("1 + 1", "3")},
			},
			PrepareCoursesRepository: addCourse,
		},
		{
			Name: "add_testing_task",
			Command: command.AddTaskCommand{
				Academic:        creator,
				CourseID:        courseID,
				TaskTitle:       "Testing task",
				TaskDescription: "Do this task",
				TaskType:        course.Testing,
				TestPoints:      []course.TestPoint{course.MustNewTestPoint("1, 2 or 3?", []string{"1", "2", "3"}, []int{0})},
			},
			PrepareCoursesRepository: addCourse,
		},
		{
			Name: "dont_add_when_invalid_task_type",
			Command: command.AddTaskCommand{
				Academic:        creator,
				CourseID:        courseID,
				TaskTitle:       "Some task",
				TaskDescription: "Don't do this task",
				TaskType:        course.TaskType(100),
			},
			PrepareCoursesRepository: addCourse,
			IsErr: func(err error) bool {
				return err != nil
			},
		},
		{
			Name: "dont_add_when_task_title_too_long",
			Command: command.AddTaskCommand{
				Academic:        creator,
				CourseID:        courseID,
				TaskTitle:       strings.Repeat("x", 201),
				TaskDescription: "Do do do",
				TaskType:        course.ManualChecking,
			},
			PrepareCoursesRepository: addCourse,
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskTitleTooLong)
			},
		},
		{
			Name: "dont_add_when_task_description_too_long",
			Command: command.AddTaskCommand{
				Academic:        creator,
				CourseID:        courseID,
				TaskTitle:       "Some course task",
				TaskDescription: strings.Repeat("x", 1001),
				TaskType:        course.Testing,
			},
			PrepareCoursesRepository: addCourse,
			IsErr: func(err error) bool {
				return errors.Is(err, course.ErrTaskDescriptionTooLong)
			},
		},
		{
			Name: "dont_add_when_academic_cant_edit_course",
			Command: command.AddTaskCommand{
				Academic:        course.MustNewAcademic("student-id", course.StudentType),
				CourseID:        courseID,
				TaskTitle:       "Some task title",
				TaskDescription: "If you want do this task",
				TaskType:        course.ManualChecking,
			},
			PrepareCoursesRepository: addCourse,
			IsErr:                    course.IsAcademicCantEditCourseError,
		},
		{
			Name: "dont_add_when_course_doesnt_exist",
			Command: command.AddTaskCommand{
				Academic:        creator,
				CourseID:        courseID,
				TaskTitle:       "Task of non-existing course",
				TaskDescription: "Don't do this",
				TaskType:        course.ManualChecking,
			},
			PrepareCoursesRepository: func(_ *course.Course) *mock.CoursesRepository {
				return mock.NewCoursesRepository()
			},
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrCourseDoesntExist)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(course.CreationParams{
				ID:      courseID,
				Creator: creator,
				Title:   "Universal course",
				Period:  course.MustNewPeriod(2043, 2044, course.FirstSemester),
			})
			coursesRepository := c.PrepareCoursesRepository(crs)
			handler := command.NewAddTaskHandler(coursesRepository)

			number, err := handler.Handle(context.Background(), c.Command)

			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				require.Equal(t, 0, number)
				return
			}
			require.NoError(t, err)
			require.Equal(t, 1, number)
			require.Equal(t, 1, crs.TasksNumber())
			task, err := crs.Task(number)
			require.NoError(t, err)
			require.Equal(t, task.Type(), c.Command.TaskType)
		})
	}
}
