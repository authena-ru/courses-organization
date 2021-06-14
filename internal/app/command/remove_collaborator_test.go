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

func TestRemoveCollaboratorHandler_Handle(t *testing.T) {
	t.Parallel()
	var (
		courseID       = "course-id"
		collaboratorID = "collaborator-id"
		creator        = course.MustNewAcademic("creator-id", course.TeacherType)
	)
	addCourse := func(crs *course.Course) *mock.CoursesRepository {
		return mock.NewCoursesRepository(crs)
	}
	testCases := []struct {
		Name                     string
		Command                  command.RemoveCollaboratorCommand
		PrepareCoursesRepository func(crs *course.Course) *mock.CoursesRepository
		IsErr                    func(err error) bool
	}{
		{
			Name: "remove_collaborator",
			Command: command.RemoveCollaboratorCommand{
				Academic:       creator,
				CourseID:       courseID,
				CollaboratorID: collaboratorID,
			},
			PrepareCoursesRepository: addCourse,
		},
		{
			Name: "dont_remove_collaborator_when_course_doesnt_exist",
			Command: command.RemoveCollaboratorCommand{
				Academic:       creator,
				CourseID:       courseID,
				CollaboratorID: collaboratorID,
			},
			PrepareCoursesRepository: func(_ *course.Course) *mock.CoursesRepository {
				return mock.NewCoursesRepository()
			},
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrCourseDoesntExist)
			},
		},
		{
			Name: "dont_remove_collaborator_when_academic_cant_edit_course",
			Command: command.RemoveCollaboratorCommand{
				Academic:       course.MustNewAcademic("other-teacher-id", course.TeacherType),
				CourseID:       courseID,
				CollaboratorID: collaboratorID,
			},
			PrepareCoursesRepository: func(crs *course.Course) *mock.CoursesRepository {
				return mock.NewCoursesRepository(crs)
			},
			IsErr: course.IsAcademicCantEditCourseError,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(course.CreationParams{
				ID:            courseID,
				Creator:       creator,
				Title:         "Chemistry",
				Period:        course.MustNewPeriod(2032, 2033, course.SecondSemester),
				Collaborators: []string{collaboratorID},
			})
			coursesRepository := c.PrepareCoursesRepository(crs)
			handler := command.NewRemoveCollaboratorHandler(coursesRepository)

			err := handler.Handle(context.Background(), c.Command)

			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			require.NotContains(t, crs.Collaborators(), c.Command.CollaboratorID)
		})
	}
}
