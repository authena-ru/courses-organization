package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

type RemoveStudentCommand struct {
	Academic  course.Academic
	CourseID  string
	StudentID string
}

type RemoveStudentHandler struct {
	coursesRepository coursesRepository
}

func NewRemoveStudentHandler(repository coursesRepository) RemoveStudentHandler {
	if repository == nil {
		panic("coursesRepository is nil")
	}
	return RemoveStudentHandler{coursesRepository: repository}
}

// Handle is RemoveStudentCommand handler.
// Removes one student from course, returns one of possible errors:
// app.ErrCourseDoesntExist, app.ErrDatabaseProblems, error that can be detected
// using method course.IsAcademicCantEditCourseError and others without definition.
func (h RemoveStudentHandler) Handle(ctx context.Context, cmd RemoveStudentCommand) error {
	err := h.coursesRepository.UpdateCourse(ctx, cmd.CourseID, removeStudent(cmd))
	return errors.Wrapf(
		err,
		"removing student #%s from course #%s by teacher #%s",
		cmd.StudentID, cmd.CourseID, cmd.Academic.ID(),
	)
}

func removeStudent(cmd RemoveStudentCommand) UpdateFunction {
	return func(_ context.Context, crs *course.Course) (*course.Course, error) {
		if err := crs.RemoveStudent(cmd.Academic, cmd.StudentID); err != nil {
			return nil, err
		}
		return crs, nil
	}
}
