package mock

import (
	"context"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/app/command"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type CoursesRepository struct {
	courses map[string]course.Course
}

func NewCoursesRepository(courses ...*course.Course) *CoursesRepository {
	crm := &CoursesRepository{
		courses: make(map[string]course.Course, len(courses)),
	}
	for _, crs := range courses {
		crm.courses[crs.ID()] = *crs
	}

	return crm
}

func (m *CoursesRepository) AddCourse(_ context.Context, crs *course.Course) error {
	m.courses[crs.ID()] = *crs

	return nil
}

func (m *CoursesRepository) GetCourse(_ context.Context, courseID string) (*course.Course, error) {
	crs, ok := m.courses[courseID]
	if !ok {
		return nil, app.ErrCourseDoesntExist
	}

	return &crs, nil
}

func (m *CoursesRepository) UpdateCourse(
	ctx context.Context,
	courseID string,
	updateFn command.UpdateFunction,
) error {
	crs, ok := m.courses[courseID]
	if !ok {
		return app.ErrCourseDoesntExist
	}

	updatedCrs, err := updateFn(ctx, &crs)
	if err != nil {
		return err
	}

	m.courses[updatedCrs.ID()] = *updatedCrs

	return nil
}

func (m *CoursesRepository) CoursesNumber() int {
	return len(m.courses)
}
