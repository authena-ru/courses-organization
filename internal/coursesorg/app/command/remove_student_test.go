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

func TestRemoveStudentHandler_Handle(t *testing.T) {
	t.Parallel()
	var (
		courseID  = "course-id"
		studentID = "student-id"
		creator   = course.MustNewAcademic("creator-id", course.Teacher)
	)
	addCourse := func(crs *course.Course) *mock.CoursesRepository {
		return mock.NewCoursesRepository(crs)
	}
	testCases := []struct {
		Name                    string
		Command                 command.RemoveStudentCommand
		PrepareCourseRepository func(crs *course.Course) *mock.CoursesRepository
		IsErr                   func(err error) bool
	}{
		{
			Name: "remove_student",
			Command: command.RemoveStudentCommand{
				Academic:  creator,
				CourseID:  courseID,
				StudentID: studentID,
			},
			PrepareCourseRepository: addCourse,
		},
		{
			Name: "dont_remove_student_when_course_doesnt_exist",
			Command: command.RemoveStudentCommand{
				Academic:  creator,
				CourseID:  courseID,
				StudentID: studentID,
			},
			PrepareCourseRepository: func(_ *course.Course) *mock.CoursesRepository {
				return mock.NewCoursesRepository()
			},
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrCourseDoesntExist)
			},
		},
		{
			Name: "dont_remove_student_when_academic_cant_edit_course",
			Command: command.RemoveStudentCommand{
				Academic:  course.MustNewAcademic("other-teacher-id", course.Teacher),
				CourseID:  courseID,
				StudentID: studentID,
			},
			PrepareCourseRepository: addCourse,
			IsErr:                   course.IsAcademicCantEditCourseError,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(course.CreationParams{
				ID:       "course-id",
				Creator:  creator,
				Title:    "Advanced Math",
				Period:   course.MustNewPeriod(2023, 2024, course.FirstSemester),
				Students: []string{studentID},
			})
			coursesRepository := c.PrepareCourseRepository(crs)
			handler := command.NewRemoveStudentHandler(coursesRepository)

			err := handler.Handle(context.Background(), c.Command)

			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			require.NotContains(t, crs.Students(), studentID)
		})
	}
}
