package query

import (
	"context"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type CoursesFilterParams struct {
	Title string
}

type allCoursesReadModel interface {
	FindAllCourses(
		ctx context.Context,
		academic course.Academic,
		filter CoursesFilterParams,
	) ([]app.CommonCourse, error)
}

type AllCoursesHandler struct {
	readModel allCoursesReadModel
}

func NewAllCoursesHandler(readModel allCoursesReadModel) AllCoursesHandler {
	if readModel == nil {
		panic("readModel is nil")
	}

	return AllCoursesHandler{readModel: readModel}
}

func (h AllCoursesHandler) Handle(ctx context.Context, qry app.AllCoursesQuery) ([]app.CommonCourse, error) {
	courses, err := h.readModel.FindAllCourses(ctx, qry.Academic, CoursesFilterParams{
		Title: qry.Title,
	})

	return courses, errors.Wrapf(err, "getting all courses of academic %v", qry.Academic)
}
