package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type ExtendCourseHandler struct {
	coursesRepository coursesRepository
}

func NewExtendCourseHandler(repository coursesRepository) ExtendCourseHandler {
	if repository == nil {
		panic("coursesRepository is nil")
	}
	return ExtendCourseHandler{coursesRepository: repository}
}

func (h ExtendCourseHandler) Handle(ctx context.Context, cmd app.ExtendCourseCommand) (extendedCourseID string, err error) {
	defer func() {
		err = errors.Wrapf(err, "extension of course #%s by teacher #%s", cmd.OriginCourseID, cmd.Academic.ID())
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

func extendCourse(extendedCourseID string, cmd app.ExtendCourseCommand) UpdateFunction {
	return func(_ context.Context, crs *course.Course) (*course.Course, error) {
		return crs.Extend(course.CreationParams{
			ID:      extendedCourseID,
			Creator: cmd.Academic,
			Title:   cmd.CourseTitle,
			Period:  cmd.CoursePeriod,
			Started: cmd.CourseStarted,
		})
	}
}
