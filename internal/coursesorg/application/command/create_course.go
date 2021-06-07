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
	academicsService  academicsService
}

func NewCreateCourseHandler(repository coursesRepository, service academicsService) CreateCourseHandler {
	if repository == nil {
		panic("coursesRepository is nil")
	}
	if service == nil {
		panic("academicsService is nil")
	}
	return CreateCourseHandler{
		coursesRepository: repository,
		academicsService:  service,
	}
}

// Handle is CreateCourseCommand handler.
// Creates course, returns id of new brand course and error.
func (h CreateCourseHandler) Handle(ctx context.Context, cmd CreateCourseCommand) (string, error) {
	courseID := uuid.NewString()
	crs, err := course.NewCourse(course.CreationParams{
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
