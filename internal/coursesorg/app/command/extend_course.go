package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

type ExtendCourseCommand struct {
	Creator        course.Academic
	OriginCourseID string
	CourseStarted  bool
	CourseTitle    string
	CoursePeriod   course.Period
}

type ExtendCourseHandler struct {
	coursesRepository coursesRepository
}

func NewExtendCourseHandler(repository coursesRepository) ExtendCourseHandler {
	if repository == nil {
		panic("coursesRepository is nil")
	}
	return ExtendCourseHandler{coursesRepository: repository}
}

// Handle is ExtendCourseCommand handler.
// Extends origin course, returns extended course ID and one of possible errors:
// app.ErrCourseDoesntExist,  course.ErrZeroCreator, course.ErrNotTeacherCantCreateCourse
// and others without definition.
func (h ExtendCourseHandler) Handle(ctx context.Context, cmd ExtendCourseCommand) (extendedCourseID string, err error) {
	defer func() {
		err = errors.Wrapf(err, "extension of course #%s by teacher #%s", cmd.OriginCourseID, cmd.Creator.ID())
	}()

	extendedCourseID = uuid.NewString()
	if err := h.coursesRepository.UpdateCourse(
		ctx,
		cmd.OriginCourseID,
		extendCourse(extendedCourseID, cmd),
	); err != nil {
		return "", err
	}
	return extendedCourseID, nil
}

func extendCourse(extendedCourseID string, cmd ExtendCourseCommand) UpdateFunction {
	return func(_ context.Context, crs *course.Course) (*course.Course, error) {
		return crs.Extend(course.CreationParams{
			ID:      extendedCourseID,
			Creator: cmd.Creator,
			Title:   cmd.CourseTitle,
			Period:  cmd.CoursePeriod,
			Started: cmd.CourseStarted,
		})
	}
}
