package command

import (
	"context"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type coursesRepository interface {
	// AddCourse returns app.ErrDatabaseProblems if repository can't
	// add course due to database problems.
	AddCourse(ctx context.Context, crs *course.Course) error

	// GetCourse returns: app.ErrCourseDoesntExist if repository can't find course,
	// app.ErrDatabaseProblems if repository can't get course due to database problems.
	GetCourse(ctx context.Context, courseID string) (*course.Course, error)

	// UpdateCourse returns: app.ErrCourseDoesntExist if repository can't find course,
	// app.ErrDatabaseProblems if repository can't update course due to database problems.
	UpdateCourse(ctx context.Context, courseID string, updateFn UpdateFunction) error
}

type UpdateFunction func(ctx context.Context, crs *course.Course) (*course.Course, error)
