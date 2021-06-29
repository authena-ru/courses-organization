package auth

import (
	"context"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type ctxKey int

const academicCtxKey ctxKey = iota

var ErrNoAcademicInContext = errors.New("no academic in context")

func AcademicFromCtx(ctx context.Context) (course.Academic, error) {
	if a, ok := ctx.Value(academicCtxKey).(course.Academic); ok {
		return a, nil
	}

	return course.Academic{}, ErrNoAcademicInContext
}
