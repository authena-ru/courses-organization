package course_test

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestNewCourse(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name        string
		Params      course.CreationParams
		ExpectedErr error
	}{
		{
			Name: "valid_course_creation_params",
			Params: course.CreationParams{
				ID:            "course-id",
				Creator:       course.MustNewAcademic("creator-id", course.TeacherType),
				Title:         "Awesome Go in backend",
				Period:        course.MustNewPeriod(2022, 2023, course.SecondSemester),
				Started:       true,
				Collaborators: []string{"collaborator-1-id", "collaborator-2-id"},
				Students:      []string{"student-id"},
			},
		},
		{
			Name: "empty_course_id",
			Params: course.CreationParams{
				Creator: course.MustNewAcademic("creator-id", course.TeacherType),
				Title:   "Programming architecture",
				Period:  course.MustNewPeriod(2021, 2022, course.FirstSemester),
				Started: false,
			},
			ExpectedErr: course.ErrEmptyCourseID,
		},
		{
			Name: "zero_course_creator",
			Params: course.CreationParams{
				ID:      "course-id",
				Title:   "JavaScript in browser",
				Period:  course.MustNewPeriod(2023, 2024, course.SecondSemester),
				Started: true,
			},
			ExpectedErr: course.ErrZeroCreator,
		},
		{
			Name: "student_cant_create_course",
			Params: course.CreationParams{
				ID:      "course-id",
				Creator: course.MustNewAcademic("student-id", course.StudentType),
				Title:   "Assembly",
				Period:  course.MustNewPeriod(2020, 2021, course.FirstSemester),
				Started: false,
			},
			ExpectedErr: course.ErrNotTeacherCantCreateCourse,
		},
		{
			Name: "empty_course_title",
			Params: course.CreationParams{
				ID:      "course-id",
				Creator: course.MustNewAcademic("creator-id", course.TeacherType),
				Period:  course.MustNewPeriod(2024, 2025, course.FirstSemester),
				Started: false,
			},
			ExpectedErr: course.ErrEmptyCourseTitle,
		},
		{
			Name: "zero_course_period",
			Params: course.CreationParams{
				ID:      "course-id",
				Creator: course.MustNewAcademic("creator-id", course.TeacherType),
				Title:   "Nice React, Awesome Angular",
				Started: true,
			},
			ExpectedErr: course.ErrZeroCoursePeriod,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs, err := course.NewCourse(c.Params)

			if c.ExpectedErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.ExpectedErr))
				return
			}
			require.NoError(t, err)
			require.Equal(t, c.Params.ID, crs.ID())
			require.Equal(t, c.Params.Creator.ID(), crs.CreatorID())
			require.Equal(t, c.Params.Title, crs.Title())
			require.Equal(t, c.Params.Period, crs.Period())
			require.Equal(t, c.Params.Started, crs.Started())
			require.ElementsMatch(t, c.Params.Students, crs.Students())
			require.ElementsMatch(t, c.Params.Collaborators, crs.Collaborators())
		})
	}
}

func TestCourse_Extend(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name              string
		Params            course.CreationParams
		NewPeriodWasGiven bool
		NewTitleWasGiven  bool
		ExpectedErr       error
	}{
		{
			Name: "extend_task_with_new_parameters",
			Params: course.CreationParams{
				ID:            "course-id",
				Creator:       course.MustNewAcademic("creator-id", course.TeacherType),
				Title:         "Clean architecture",
				Period:        course.MustNewPeriod(2027, 2028, course.FirstSemester),
				Students:      []string{"some-student-id"},
				Collaborators: []string{"some-collaborator-id"},
			},
			NewPeriodWasGiven: true,
			NewTitleWasGiven:  true,
		},
		{
			Name: "extend_task_without_new_period_gives_new_task_with_next_period",
			Params: course.CreationParams{
				ID:      "course-id",
				Creator: course.MustNewAcademic("creator-id", course.TeacherType),
				Title:   "Clean Clean Clean",
			},
			NewTitleWasGiven: true,
		},
		{
			Name: "extend_task_without_new_title_gives_new_task_with_origin_title",
			Params: course.CreationParams{
				ID:      "course-id",
				Creator: course.MustNewAcademic("creator-id", course.TeacherType),
				Period:  course.MustNewPeriod(2030, 2031, course.FirstSemester),
			},
			NewPeriodWasGiven: true,
		},
		{
			Name: "empty_course_id",
			Params: course.CreationParams{
				Creator: course.MustNewAcademic("teacher-id", course.TeacherType),
			},
			ExpectedErr: course.ErrEmptyCourseID,
		},
		{
			Name: "zero_creator",
			Params: course.CreationParams{
				ID: "course-id",
			},
			ExpectedErr: course.ErrZeroCreator,
		},
		{
			Name: "not_teacher_cant_extend_course",
			Params: course.CreationParams{
				ID:      "course-id",
				Creator: course.MustNewAcademic("student-id", course.StudentType),
			},
			ExpectedErr: course.ErrNotTeacherCantCreateCourse,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			creator := course.MustNewAcademic("origin-course-creator-id", course.TeacherType)
			originCourse := course.MustNewCourse(course.CreationParams{
				ID:            "origin-course-id",
				Creator:       creator,
				Title:         "Architecture",
				Period:        course.MustNewPeriod(2024, 2025, course.FirstSemester),
				Students:      []string{"student-id"},
				Collaborators: []string{"collaborator-id¬"},
			})
			_, err := originCourse.AddManualCheckingTask(creator, course.ManualCheckingTaskCreationParams{
				Title:       "Adapters",
				Description: "Write your adapters",
				Deadline: course.MustNewDeadline(
					time.Date(2025, time.September, 1, 0, 0, 0, 0, time.Local),
					time.Date(2025, time.September, 15, 0, 0, 0, 0, time.Local),
				),
			})
			require.NoError(t, err)
			_, err = originCourse.AddAutoCodeCheckingTask(creator, course.AutoCodeCheckingTaskCreationParams{
				Title:       "Printer class",
				Description: "Write your Printer",
				Deadline: course.MustNewDeadline(
					time.Date(2025, time.October, 1, 0, 0, 0, 0, time.Local),
					time.Date(2025, time.October, 17, 0, 0, 0, 0, time.Local),
				),
				TestData: []course.TestData{course.MustNewTestData("1", "Print: 1")},
			})
			require.NoError(t, err)
			_, err = originCourse.AddTestingTask(creator, course.TestingTaskCreationParams{
				Title:       "Entities",
				Description: "Entities test",
				TestPoints:  []course.TestPoint{course.MustNewTestPoint("Entities are classes", []string{"Yes", "No"}, []int{1})},
			})
			require.NoError(t, err)

			extendedCourse, err := originCourse.Extend(c.Params)

			if c.ExpectedErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.ExpectedErr))
				return
			}
			require.NoError(t, err)
			require.Equal(t, c.Params.ID, extendedCourse.ID())
			require.Equal(t, c.Params.Creator.ID(), extendedCourse.CreatorID())
			if c.NewPeriodWasGiven {
				require.Equal(t, c.Params.Period, extendedCourse.Period())
			} else {
				require.Equal(t, course.MustNewPeriod(2025, 2026, course.FirstSemester), extendedCourse.Period())
			}
			if c.NewTitleWasGiven {
				require.Equal(t, c.Params.Title, extendedCourse.Title())
			} else {
				require.Equal(t, originCourse.Title(), extendedCourse.Title())
			}
			require.ElementsMatch(t, append(originCourse.Students(), c.Params.Students...), extendedCourse.Students())
			require.ElementsMatch(t, append(originCourse.Collaborators(), c.Params.Collaborators...), extendedCourse.Collaborators())
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
		})
	}
}
