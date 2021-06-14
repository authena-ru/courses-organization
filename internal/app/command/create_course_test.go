package command_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/app/command"
	"github.com/authena-ru/courses-organization/internal/app/command/mock"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestCreateCourseHandler_Handle(t *testing.T) {
	t.Parallel()
	creator := course.MustNewAcademic("creator-id", course.Teacher)
	testCases := []struct {
		Name        string
		Command     command.CreateCourseCommand
		ExpectedErr error
	}{
		{
			Name: "create_course",
			Command: command.CreateCourseCommand{
				Academic:      creator,
				CourseStarted: true,
				CourseTitle:   "Bla Bla Literature",
				CoursePeriod:  course.MustNewPeriod(2019, 2020, course.FirstSemester),
			},
		},
		{
			Name: "dont_create_when_zero_creator",
			Command: command.CreateCourseCommand{
				CourseStarted: false,
				CourseTitle:   "Bla Literature",
				CoursePeriod:  course.MustNewPeriod(2024, 2025, course.SecondSemester),
			},
			ExpectedErr: course.ErrZeroCreator,
		},
		{
			Name: "dont_create_when_empty_course_title",
			Command: command.CreateCourseCommand{
				Academic:      creator,
				CourseStarted: true,
				CoursePeriod:  course.MustNewPeriod(2040, 2041, course.FirstSemester),
			},
			ExpectedErr: course.ErrEmptyCourseTitle,
		},
		{
			Name: "dont_create_when_zero_course_period",
			Command: command.CreateCourseCommand{
				Academic:      creator,
				CourseStarted: false,
				CourseTitle:   "Literature",
			},
			ExpectedErr: course.ErrZeroCoursePeriod,
		},
		{
			Name: "dont_create_when_not_teacher_creates_course",
			Command: command.CreateCourseCommand{
				Academic:      course.MustNewAcademic("student-id", course.Student),
				CourseStarted: false,
				CourseTitle:   "Literature bla",
				CoursePeriod:  course.MustNewPeriod(2024, 2025, course.SecondSemester),
			},
			ExpectedErr: course.ErrNotTeacherCantCreateCourse,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			coursesRepository := mock.NewCoursesRepository()
			handler := command.NewCreateCourseHandler(coursesRepository)

			courseID, err := handler.Handle(context.Background(), c.Command)

			if c.ExpectedErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.ExpectedErr))
				require.Empty(t, courseID)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, courseID)
			require.Equal(t, 1, coursesRepository.CoursesNumber())
		})
	}
}
