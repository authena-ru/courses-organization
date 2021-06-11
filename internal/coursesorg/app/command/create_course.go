package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

type CreateCourseCommand struct {
	Academic      course.Academic
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
// Creates course, returns ID of new brand course and one of possible errors:
// course.ErrZeroCreator, course.ErrNotTeacherCantCreateCourse, course.ErrEmptyCourseTitle,
// course.ErrZeroCoursePeriod and others without definition.
func (h CreateCourseHandler) Handle(ctx context.Context, cmd CreateCourseCommand) (courseID string, err error) {
	defer func() {
		err = errors.Wrapf(err, "course creation by teacher #%s", cmd.Academic.ID())
	}()

	courseID = uuid.NewString()
	crs, err := course.NewCourse(course.CreationParams{
		ID:      courseID,
		Creator: cmd.Academic,
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
	return
}
