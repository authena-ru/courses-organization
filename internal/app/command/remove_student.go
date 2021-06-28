package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type RemoveStudentHandler struct {
	coursesRepository coursesRepository
}

func NewRemoveStudentHandler(repository coursesRepository) RemoveStudentHandler {
	if repository == nil {
		panic("coursesRepository is nil")
	}
	return RemoveStudentHandler{coursesRepository: repository}
}

func (h RemoveStudentHandler) Handle(ctx context.Context, cmd app.RemoveStudentCommand) error {
	err := h.coursesRepository.UpdateCourse(ctx, cmd.CourseID, removeStudent(cmd))
	return errors.Wrapf(
		err,
		"removing student #%s from course #%s by teacher #%s",
		cmd.StudentID, cmd.CourseID, cmd.Academic.ID(),
	)
}

func removeStudent(cmd app.RemoveStudentCommand) UpdateFunction {
	return func(_ context.Context, crs *course.Course) (*course.Course, error) {
		if err := crs.RemoveStudent(cmd.Academic, cmd.StudentID); err != nil {
			return nil, err
		}
		return crs, nil
	}
}
