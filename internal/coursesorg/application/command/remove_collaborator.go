package command

import (
	"context"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

type RemoveCollaboratorCommand struct {
	Teacher        course.Academic
	CourseID       string
	CollaboratorID string
}

type RemoveCollaboratorHandler struct {
	coursesRepository coursesRepository
	academicsService  academicsService
}

func NewRemoveCollaboratorHandler(repository coursesRepository, service academicsService) RemoveCollaboratorHandler {
	if repository == nil {
		panic("coursesRepository is nil")
	}
	if service == nil {
		panic("academicsService is nil")
	}
	return RemoveCollaboratorHandler{
		coursesRepository: repository,
		academicsService:  service,
	}
}

// Handle is RemoveCollaboratorCommand handler.
// Removes one collaborator from course, returns error.
func (h RemoveCollaboratorHandler) Handle(ctx context.Context, cmd RemoveCollaboratorCommand) error {
	if err := h.academicsService.TeacherExists(cmd.Teacher.ID()); err != nil {
		return err
	}
	return h.coursesRepository.UpdateCourse(ctx, cmd.CourseID, cmd.Teacher, removeCollaborator(cmd))
}

func removeCollaborator(cmd RemoveCollaboratorCommand) UpdateFunction {
	return func(_ context.Context, crs *course.Course) (*course.Course, error) {
		if err := crs.RemoveCollaborator(cmd.Teacher, cmd.CollaboratorID); err != nil {
			return nil, err
		}
		return crs, nil
	}
}
