package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

type AddCollaboratorCommand struct {
	Teacher        course.Academic
	CourseID       string
	CollaboratorID string
}

type AddCollaboratorHandler struct {
	coursesRepository coursesRepository
	academicsService  academicsService
}

func NewAddCollaboratorHandler(repository coursesRepository, service academicsService) AddCollaboratorHandler {
	if repository == nil {
		panic("coursesRepository is nil")
	}
	if service == nil {
		panic("academicsService is nil")
	}
	return AddCollaboratorHandler{
		coursesRepository: repository,
		academicsService:  service,
	}
}

// Handle is AddCollaboratorCommand handler.
// Adds one collaborator to course, returns one of possible errors:
// app.ErrTeacherDoesntExist, app.ErrCourseDoesntExist, error that
// can be detected using course.IsAcademicCantEditCourseError and
// others without definition.
func (h AddCollaboratorHandler) Handle(ctx context.Context, cmd AddCollaboratorCommand) (err error) {
	defer func() {
		err = errors.Wrapf(
			err,
			"adding collaborator #%s to course #%s by teacher #%s",
			cmd.CollaboratorID, cmd.CourseID, cmd.Teacher.ID(),
		)
	}()

	if err := h.academicsService.TeacherExists(cmd.CollaboratorID); err != nil {
		return err
	}
	return h.coursesRepository.UpdateCourse(ctx, cmd.CourseID, addCollaborator(cmd))
}

func addCollaborator(cmd AddCollaboratorCommand) UpdateFunction {
	return func(_ context.Context, crs *course.Course) (*course.Course, error) {
		if err := crs.AddCollaborators(cmd.Teacher, cmd.CollaboratorID); err != nil {
			return nil, err
		}
		return crs, nil
	}
}
