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

func TestAddCollaboratorHandler_Handle(t *testing.T) {
	t.Parallel()

	addCourse := func(crs *course.Course) *mock.CoursesRepository {
		return mock.NewCoursesRepository(crs)
	}
	addCollaborator := func() *mock.AcademicsService {
		return mock.NewAcademicsService([]string{"collaborator-id"}, nil, nil)
	}
	testCases := []struct {
		Name                     string
		Command                  app.AddCollaboratorCommand
		PrepareCoursesRepository func(crs *course.Course) *mock.CoursesRepository
		PrepareAcademicsService  func() *mock.AcademicsService
		IsErr                    func(err error) bool
	}{
		{
			Name: "add_collaborator",
			Command: app.AddCollaboratorCommand{
				Academic:       course.MustNewAcademic("creator-id", course.TeacherType),
				CourseID:       "course-id",
				CollaboratorID: "collaborator-id",
			},
			PrepareCoursesRepository: addCourse,
			PrepareAcademicsService:  addCollaborator,
		},
		{
			Name: "dont_add_when_teacher_cant_edit_course",
			Command: app.AddCollaboratorCommand{
				Academic:       course.MustNewAcademic("other-creator-id", course.TeacherType),
				CourseID:       "course-id",
				CollaboratorID: "collaborator-id",
			},
			PrepareCoursesRepository: addCourse,
			PrepareAcademicsService:  addCollaborator,
			IsErr:                    course.IsAcademicCantEditCourseError,
		},
		{
			Name: "dont_add_when_collaborator_doesnt_exist_as_teacher",
			Command: app.AddCollaboratorCommand{
				Academic:       course.MustNewAcademic("creator-id", course.TeacherType),
				CourseID:       "course-id",
				CollaboratorID: "collaborator-id",
			},
			PrepareCoursesRepository: addCourse,
			PrepareAcademicsService: func() *mock.AcademicsService {
				return mock.NewAcademicsService(nil, nil, nil)
			},
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrTeacherDoesntExist)
			},
		},
		{
			Name: "dont_add_when_course_doesnt_exist",
			Command: app.AddCollaboratorCommand{
				Academic:       course.MustNewAcademic("creator-id", course.TeacherType),
				CourseID:       "course-id",
				CollaboratorID: "collaborator-id",
			},
			PrepareCoursesRepository: func(_ *course.Course) *mock.CoursesRepository {
				return mock.NewCoursesRepository()
			},
			PrepareAcademicsService: addCollaborator,
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
				ID:      "course-id",
				Creator: course.MustNewAcademic("creator-id", course.TeacherType),
				Title:   "Docker and Kubernetes",
				Period:  course.MustNewPeriod(2023, 2024, course.FirstSemester),
			})
			coursesRepository := c.PrepareCoursesRepository(crs)
			academicsService := c.PrepareAcademicsService()
			handler := command.NewAddCollaboratorHandler(coursesRepository, academicsService)

			err := handler.Handle(context.Background(), c.Command)

			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				require.NotContains(t, crs.Collaborators(), c.Command.CollaboratorID)

				return
			}
			require.NoError(t, err)
			require.Contains(t, crs.Collaborators(), c.Command.CollaboratorID)
		})
	}
}
