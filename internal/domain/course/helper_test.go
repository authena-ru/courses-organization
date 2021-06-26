package course_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type CourseOption func(academic course.Academic, crs *course.Course)

func WithStudents(students ...string) CourseOption {
	return func(academic course.Academic, crs *course.Course) {
		err := crs.AddStudents(academic, students...)
		if err != nil {
			panic(err)
		}
	}
}

func WithCollaborators(collaborators ...string) CourseOption {
	return func(academic course.Academic, crs *course.Course) {
		err := crs.AddCollaborators(academic, collaborators...)
		if err != nil {
			panic(err)
		}
	}
}

func NewCourse(t *testing.T, creator course.Academic, opts ...CourseOption) *course.Course {
	t.Helper()
	crs := course.MustNewCourse(course.CreationParams{
		ID:      "origin-course-id",
		Creator: creator,
		Title:   "Course title",
		Period:  course.MustNewPeriod(2024, 2025, course.FirstSemester),
	})
	for _, opt := range opts {
		opt(creator, crs)
	}
	return crs
}

func AddManualCheckingTaskToCourse(t *testing.T, academic course.Academic, crs *course.Course) int {
	taskNumber, err := crs.AddManualCheckingTask(academic, course.ManualCheckingTaskCreationParams{
		Title:       "Manual checking task title",
		Description: "Manual checking task description",
		Deadline: course.MustNewDeadline(
			time.Date(2025, time.September, 1, 0, 0, 0, 0, time.Local),
			time.Date(2025, time.September, 15, 0, 0, 0, 0, time.Local),
		),
	})
	require.NoError(t, err)
	return taskNumber
}

func AddAutoCodeCheckingTaskToCourse(t *testing.T, academic course.Academic, crs *course.Course) int {
	taskNumber, err := crs.AddAutoCodeCheckingTask(academic, course.AutoCodeCheckingTaskCreationParams{
		Title:       "Auto code checking task title",
		Description: "Auto code checking task description",
		Deadline: course.MustNewDeadline(
			time.Date(2025, time.October, 1, 0, 0, 0, 0, time.Local),
			time.Date(2025, time.October, 17, 0, 0, 0, 0, time.Local),
		),
		TestData: []course.TestData{course.MustNewTestData("1", "Print: 1")},
	})
	require.NoError(t, err)
	return taskNumber
}

func AddTestingTaskToCourse(t *testing.T, academic course.Academic, crs *course.Course) int {
	taskNumber, err := crs.AddTestingTask(academic, course.TestingTaskCreationParams{
		Title:       "Testing task title",
		Description: "Testing task description",
		TestPoints:  []course.TestPoint{course.MustNewTestPoint("Yes/no question", []string{"Yes", "No"}, []int{1})},
	})
	require.NoError(t, err)
	return taskNumber
}

func requireExtendedCourse(
	t *testing.T,
	originCourse *course.Course,
	extendedCourse *course.Course,
	params course.CreationParams,
	newTitleWasGiven, newPeriodWasGiven bool,
) {
	require.Equal(t, params.ID, extendedCourse.ID())
	require.Equal(t, params.Creator.ID(), extendedCourse.CreatorID())
	if newPeriodWasGiven {
		require.Equal(t, params.Period, extendedCourse.Period())
	} else {
		require.Equal(t, course.MustNewPeriod(2025, 2026, course.FirstSemester), extendedCourse.Period())
	}
	if newTitleWasGiven {
		require.Equal(t, params.Title, extendedCourse.Title())
	} else {
		require.Equal(t, originCourse.Title(), extendedCourse.Title())
	}
	require.ElementsMatch(t, append(originCourse.Students(), params.Students...), extendedCourse.Students())
	require.ElementsMatch(t, append(originCourse.Collaborators(), params.Collaborators...), extendedCourse.Collaborators())
	requireCourseTasksEquals(t, originCourse, extendedCourse)
}

func requireCourseTasksEquals(t *testing.T, originCourse, extendedCourse *course.Course) {
	t.Helper()
	require.Equal(t, originCourse.TasksNumber(), extendedCourse.TasksNumber())
	for i := 1; i <= extendedCourse.TasksNumber(); i++ {
		taskFromOrigin, err := originCourse.Task(i)
		require.NoError(t, err)
		taskFromExtended, err := extendedCourse.Task(i)
		require.NoError(t, err)
		require.Equal(t, taskFromOrigin.Number(), taskFromExtended.Number())
		require.Equal(t, taskFromOrigin.Title(), taskFromExtended.Title())
		require.Equal(t, taskFromOrigin.Description(), taskFromExtended.Description())
		require.Equal(t, taskFromOrigin.Type(), taskFromExtended.Type())
		extendedDeadline, _ := taskFromExtended.Deadline()
		require.True(t, extendedDeadline.IsZero())
		originTestData, _ := taskFromOrigin.TestData()
		extendedTestData, _ := taskFromExtended.TestData()
		require.Equal(t, originTestData, extendedTestData)
		originTestPoints, _ := taskFromOrigin.TestPoints()
		extendedTestPoints, _ := taskFromExtended.TestPoints()
		require.Equal(t, originTestPoints, extendedTestPoints)
	}
}

func requireGeneralTaskParamsEquals(
	t *testing.T,
	task course.Task,
	number int, taskType course.TaskType, title, description string,
) {
	require.Equal(t, number, task.Number())
	require.Equal(t, taskType, task.Type())
	require.Equal(t, title, task.Title())
	require.Equal(t, description, task.Description())
}
