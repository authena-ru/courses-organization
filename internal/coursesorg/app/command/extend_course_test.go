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
	addOriginCourse := func(crs *course.Course, crm *mock.CoursesRepository) {
		crm.Courses = map[string]course.Course{crs.ID(): *crs}
	}
	testCases := []struct {
		Name                     string
		Command                  command.ExtendCourseCommand
		PrepareCoursesRepository func(crs *course.Course, crm *mock.CoursesRepository)
		ExpectedErr              error
	}{
		{
			Name: "add_when_origin_course_exists",
			Command: command.ExtendCourseCommand{
				Creator:        creator,
				OriginCourseID: originCourseID,
				CourseStarted:  false,
				CourseTitle:    "Phy Physics",
				CoursePeriod:   course.MustNewPeriod(2026, 2027, course.SecondSemester),
			},
			PrepareCoursesRepository: addOriginCourse,
		},
		{
			Name: "dont_add_when_origin_course_doesnt_exist",
			Command: command.ExtendCourseCommand{
				Creator:        creator,
				OriginCourseID: originCourseID,
				CourseStarted:  true,
			},
			PrepareCoursesRepository: func(_ *course.Course, crm *mock.CoursesRepository) {
				crm.Courses = make(map[string]course.Course)
			},
			ExpectedErr: app.ErrCourseDoesntExist,
		},
		{
			Name: "dont_add_when_zero_creator",
			Command: command.ExtendCourseCommand{
				OriginCourseID: originCourseID,
				CourseStarted:  false,
			},
			PrepareCoursesRepository: addOriginCourse,
			ExpectedErr:              course.ErrZeroCreator,
		},
		{
			Name: "dont_add_when_not_teacher_extends_course",
			Command: command.ExtendCourseCommand{
				OriginCourseID: originCourseID,
				Creator:        course.MustNewAcademic("student-id", course.Student),
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
			coursesRepository := &mock.CoursesRepository{}
			c.PrepareCoursesRepository(originCourse, coursesRepository)
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
		})
	}
}
