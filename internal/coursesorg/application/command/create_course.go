package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

type CreateCourseCommand struct {
	Creator       course.Academic
	CourseStarted bool
	CourseTitle   string
	CoursePeriod  course.Period
}

type CreateCourseHandler struct {
	coursesRepository coursesRepository
}

func NewCreateCourseHandler(repository coursesRepository) CreateCourseHandler {
	if repository == nil {
		panic("coursesRepository is nil")
	}
	return CreateCourseHandler{coursesRepository: repository}
}

// Handle is CreateCourseCommand handler.
// Creates course, returns id of new brand course and error.
func (h CreateCourseHandler) Handle(ctx context.Context, cmd CreateCourseCommand) (string, error) {
	courseID := uuid.NewString()
	crs, err := course.NewCourse(course.CreationCourseParams{
		ID:      courseID,
		Creator: cmd.Creator,
		Title:   cmd.CourseTitle,
		Period:  cmd.CoursePeriod,
		Started: cmd.CourseStarted,
	})
	if err != nil {
		return "", err
	}
	if err := h.coursesRepository.AddCourse(ctx, crs); err != nil {
		return "", err
	}
	return courseID, nil
}
