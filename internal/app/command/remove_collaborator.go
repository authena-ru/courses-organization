package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type RemoveCollaboratorHandler struct {
	coursesRepository coursesRepository
}

func NewRemoveCollaboratorHandler(repository coursesRepository) RemoveCollaboratorHandler {
	if repository == nil {
		panic("coursesRepository is nil")
	}
	return RemoveCollaboratorHandler{coursesRepository: repository}
}

func (h RemoveCollaboratorHandler) Handle(ctx context.Context, cmd app.RemoveCollaboratorCommand) error {
	err := h.coursesRepository.UpdateCourse(ctx, cmd.CourseID, removeCollaborator(cmd))
	return errors.Wrapf(
		err,
		"removing collaborator #%s from course #%s by teacher #%s",
		cmd.CollaboratorID, cmd.CollaboratorID, cmd.Academic.ID(),
	)
}

func removeCollaborator(cmd app.RemoveCollaboratorCommand) UpdateFunction {
	return func(_ context.Context, crs *course.Course) (*course.Course, error) {
		if err := crs.RemoveCollaborator(cmd.Academic, cmd.CollaboratorID); err != nil {
			return nil, err
		}
		return crs, nil
	}
}
