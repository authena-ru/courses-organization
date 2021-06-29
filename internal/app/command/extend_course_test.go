package command_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/app/command"
	"github.com/authena-ru/courses-organization/internal/app/command/mock"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestExtendCourseHandler_Handle(t *testing.T) {
	t.Parallel()

	addOriginCourse := func(crs *course.Course) *mock.CoursesRepository {
		return mock.NewCoursesRepository(crs)
	}
	testCases := []struct {
		Name                     string
		Command                  app.ExtendCourseCommand
		PrepareCoursesRepository func(crs *course.Course) *mock.CoursesRepository
		ExpectedErr              error
	}{
		{
			Name: "extend_existing_origin_course",
			Command: app.ExtendCourseCommand{
				Academic:       course.MustNewAcademic("creator-id", course.TeacherType),
				OriginCourseID: "origin-course-id",
				CourseStarted:  false,
				CourseTitle:    "Phy Physics",
				CoursePeriod:   course.MustNewPeriod(2026, 2027, course.SecondSemester),
			},
			PrepareCoursesRepository: addOriginCourse,
		},
		{
			Name: "dont_extend_when_origin_course_doesnt_exist",
			Command: app.ExtendCourseCommand{
				Academic:       course.MustNewAcademic("creator-id", course.TeacherType),
				OriginCourseID: "origin-course-id",
				CourseStarted:  true,
			},
			PrepareCoursesRepository: func(_ *course.Course) *mock.CoursesRepository {
				return mock.NewCoursesRepository()
			},
			ExpectedErr: app.ErrCourseDoesntExist,
		},
		{
			Name: "dont_extend_when_zero_creator",
			Command: app.ExtendCourseCommand{
				OriginCourseID: "origin-course-id",
				CourseStarted:  false,
			},
			PrepareCoursesRepository: addOriginCourse,
			ExpectedErr:              course.ErrZeroCreator,
		},
		{
			Name: "dont_extend_when_not_teacher_extends_course",
			Command: app.ExtendCourseCommand{
				OriginCourseID: "origin-course-id",
				Academic:       course.MustNewAcademic("student-id", course.StudentType),
				CourseStarted:  true,
			},
			PrepareCoursesRepository: addOriginCourse,
			ExpectedErr:              course.ErrNotTeacherCantCreateCourse,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			originCourse := course.MustNewCourse(course.CreationParams{
				ID:      "origin-course-id",
				Creator: course.MustNewAcademic("creator-id", course.TeacherType),
				Title:   "Physics",
				Period:  course.MustNewPeriod(2023, 2024, course.FirstSemester),
			})
			coursesRepository := c.PrepareCoursesRepository(originCourse)
			handler := command.NewExtendCourseHandler(coursesRepository)

			extendedCourseID, err := handler.Handle(context.Background(), c.Command)

			if c.ExpectedErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.ExpectedErr))
				require.Empty(t, extendedCourseID)

				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, extendedCourseID)
			require.Equal(t, 2, coursesRepository.CoursesNumber())
		})
	}
}
