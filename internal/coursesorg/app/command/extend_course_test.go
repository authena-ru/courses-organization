package command_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/coursesorg/app"
	"github.com/authena-ru/courses-organization/internal/coursesorg/app/command"
	"github.com/authena-ru/courses-organization/internal/coursesorg/app/command/mock"
	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

func TestExtendCourseHandler_Handle(t *testing.T) {
	t.Parallel()
	var (
		originCourseID = "origin-course-id"
		creator        = course.MustNewAcademic("creator-id", course.Teacher)
	)
	addOriginCourse := func(crs *course.Course) *mock.CoursesRepository {
		return mock.NewCoursesRepository(crs)
	}
	testCases := []struct {
		Name                     string
		Command                  command.ExtendCourseCommand
		PrepareCoursesRepository func(crs *course.Course) *mock.CoursesRepository
		ExpectedErr              error
	}{
		{
			Name: "extend_existing_origin_course",
			Command: command.ExtendCourseCommand{
				Academic:       creator,
				OriginCourseID: originCourseID,
				CourseStarted:  false,
				CourseTitle:    "Phy Physics",
				CoursePeriod:   course.MustNewPeriod(2026, 2027, course.SecondSemester),
			},
			PrepareCoursesRepository: addOriginCourse,
		},
		{
			Name: "dont_extend_when_origin_course_doesnt_exist",
			Command: command.ExtendCourseCommand{
				Academic:       creator,
				OriginCourseID: originCourseID,
				CourseStarted:  true,
			},
			PrepareCoursesRepository: func(_ *course.Course) *mock.CoursesRepository {
				return mock.NewCoursesRepository()
			},
			ExpectedErr: app.ErrCourseDoesntExist,
		},
		{
			Name: "dont_extend_when_zero_creator",
			Command: command.ExtendCourseCommand{
				OriginCourseID: originCourseID,
				CourseStarted:  false,
			},
			PrepareCoursesRepository: addOriginCourse,
			ExpectedErr:              course.ErrZeroCreator,
		},
		{
			Name: "dont_extend_when_not_teacher_extends_course",
			Command: command.ExtendCourseCommand{
				OriginCourseID: originCourseID,
				Academic:       course.MustNewAcademic("student-id", course.Student),
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
				ID:      originCourseID,
				Creator: creator,
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
