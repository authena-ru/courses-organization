package query

import (
	"context"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type specificCourseReadModel interface {
	FindCourse(ctx context.Context, academic course.Academic, courseID string) (app.CommonCourse, error)
}

type SpecificCourseHandler struct {
	readModel specificCourseReadModel
}

func NewSpecificCourseHandler(readModel specificCourseReadModel) SpecificCourseHandler {
	if readModel == nil {
		panic("readModel is nil")
	}

	return SpecificCourseHandler{readModel: readModel}
}

func (h SpecificCourseHandler) Handle(ctx context.Context, qry app.SpecificCourseQuery) (app.CommonCourse, error) {
	crs, err := h.readModel.FindCourse(ctx, qry.Academic, qry.CourseID)

	return crs, errors.Wrapf(err, "getting course #%s by academic %v", qry.CourseID, qry.Academic)
}
