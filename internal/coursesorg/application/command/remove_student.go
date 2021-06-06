package command

import (
	"context"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

type RemoveStudentCommand struct {
	Teacher   course.Academic
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
// Removes one student from course, returns error.
func (h RemoveStudentHandler) Handle(ctx context.Context, cmd RemoveStudentCommand) error {
	return h.coursesRepository.UpdateCourse(ctx, cmd.CourseID, cmd.Teacher, removeStudent(cmd))
}

func removeStudent(cmd RemoveStudentCommand) UpdateFunction {
	return func(_ context.Context, crs *course.Course) (*course.Course, error) {
		if err := crs.RemoveStudent(cmd.Teacher, cmd.StudentID); err != nil {
			return nil, err
		}
		return crs, nil
	}
}
