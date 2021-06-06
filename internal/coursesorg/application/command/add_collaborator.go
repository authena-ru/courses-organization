package command

import (
	"context"

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
		panic("CoursesRepository is nil")
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
// Adds one collaborator to course, returns error.
func (h AddCollaboratorHandler) Handle(ctx context.Context, cmd AddCollaboratorCommand) error {
	if err := h.academicsService.TeacherExists(cmd.CollaboratorID); err != nil {
		return err
	}
	return h.coursesRepository.UpdateCourse(ctx, cmd.CourseID, cmd.Teacher, addCollaborator(cmd))
}

func addCollaborator(cmd AddCollaboratorCommand) UpdateFunction {
	return func(_ context.Context, crs *course.Course) (*course.Course, error) {
		if err := crs.AddCollaborators(cmd.Teacher, cmd.CollaboratorID); err != nil {
			return nil, err
		}
		return crs, nil
	}
}
