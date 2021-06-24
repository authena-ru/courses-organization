package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type AddCollaboratorCommand struct {
	Academic       course.Academic
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
// app.ErrTeacherDoesntExist, app.ErrCourseDoesntExist, app.ErrDatabaseProblems,
// error that can be detected using course.IsAcademicCantEditCourseError and
// others without definition.
func (h AddCollaboratorHandler) Handle(ctx context.Context, cmd AddCollaboratorCommand) error {
	err := h.coursesRepository.UpdateCourse(ctx, cmd.CourseID, h.addCollaborator(cmd))
	return errors.Wrapf(
		err,
		"adding collaborator #%s to course #%s by academic #%s",
		cmd.CollaboratorID, cmd.CourseID, cmd.Academic.ID(),
	)
}

func (h AddCollaboratorHandler) addCollaborator(cmd AddCollaboratorCommand) UpdateFunction {
	return func(ctx context.Context, crs *course.Course) (*course.Course, error) {
		if err := h.academicsService.TeacherExists(ctx, cmd.CollaboratorID); err != nil {
			return nil, err
		}
		if err := crs.AddCollaborators(cmd.Academic, cmd.CollaboratorID); err != nil {
			return nil, err
		}
		return crs, nil
	}
}
