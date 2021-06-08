package command

import (
	"context"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

type coursesRepository interface {
	AddCourse(ctx context.Context, crs *course.Course) error

	GetCourse(ctx context.Context, courseID string) (*course.Course, error)

	UpdateCourse(ctx context.Context, courseID string, updateFn UpdateFunction) error
}

type UpdateFunction func(ctx context.Context, crs *course.Course) (*course.Course, error)
