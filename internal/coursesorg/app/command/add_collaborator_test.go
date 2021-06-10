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

func TestAddCollaboratorHandler_Handle(t *testing.T) {
	t.Parallel()
	var (
		courseID       = "course-id"
		collaboratorID = "collaborator-id"
		creator        = course.MustNewAcademic("creator-id", course.Teacher)
	)
	addCourse := func(crs *course.Course, crm *mock.CoursesRepository) {
		crm.Courses = map[string]course.Course{crs.ID(): *crs}
	}
	addCollaborator := func(asm *mock.AcademicsService) {
		asm.Teachers = map[string]bool{collaboratorID: true}
	}
	testCases := []struct {
		Name                     string
		Command                  command.AddCollaboratorCommand
		PrepareCoursesRepository func(crs *course.Course, crm *mock.CoursesRepository)
		PrepareAcademicsService  func(asm *mock.AcademicsService)
		IsErr                    func(err error) bool
	}{
		{
			Name: "add_collaborator",
			Command: command.AddCollaboratorCommand{
				Teacher:        creator,
				CourseID:       courseID,
				CollaboratorID: collaboratorID,
			},
			PrepareCoursesRepository: addCourse,
			PrepareAcademicsService:  addCollaborator,
		},
		{
			Name: "dont_add_when_teacher_cant_edit_course",
			Command: command.AddCollaboratorCommand{
				Teacher:        course.MustNewAcademic("other-creator-id", course.Teacher),
				CourseID:       courseID,
				CollaboratorID: collaboratorID,
			},
			PrepareCoursesRepository: addCourse,
			PrepareAcademicsService:  addCollaborator,
			IsErr:                    course.IsAcademicCantEditCourseError,
		},
		{
			Name: "dont_add_when_collaborator_doesnt_exist_as_teacher",
			Command: command.AddCollaboratorCommand{
				Teacher:        creator,
				CourseID:       courseID,
				CollaboratorID: collaboratorID,
			},
			PrepareCoursesRepository: addCourse,
			PrepareAcademicsService: func(asm *mock.AcademicsService) {
				asm.Teachers = make(map[string]bool)
			},
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrTeacherDoesntExist)
			},
		},
		{
			Name: "dont_add_when_update_fails",
			Command: command.AddCollaboratorCommand{
				Teacher:        creator,
				CourseID:       courseID,
				CollaboratorID: collaboratorID,
			},
			PrepareCoursesRepository: func(_ *course.Course, crm *mock.CoursesRepository) {
				crm.Courses = make(map[string]course.Course)
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
				ID:      courseID,
				Creator: creator,
				Title:   "Docker and Kubernetes",
				Period:  course.MustNewPeriod(2023, 2024, course.FirstSemester),
			})
			coursesRepository := &mock.CoursesRepository{}
			c.PrepareCoursesRepository(crs, coursesRepository)
			academicsService := &mock.AcademicsService{}
			c.PrepareAcademicsService(academicsService)
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
