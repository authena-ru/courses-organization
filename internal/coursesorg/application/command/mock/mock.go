package mock

import (
	"context"

	"github.com/authena-ru/courses-organization/internal/coursesorg/application/apperr"
	"github.com/authena-ru/courses-organization/internal/coursesorg/application/command"
	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

type CoursesRepository struct {
	Courses map[string]course.Course
}

func (m *CoursesRepository) AddCourse(_ context.Context, crs *course.Course) error {
	m.Courses[crs.ID()] = *crs
	return nil
}

func (m *CoursesRepository) GetCourse(_ context.Context, courseID string) (*course.Course, error) {
	crs, ok := m.Courses[courseID]
	if !ok {
		return nil, apperr.ErrCourseNotFound
	}
	return &crs, nil
}

func (m *CoursesRepository) UpdateCourse(
	ctx context.Context,
	courseID string,
	updateFn command.UpdateFunction,
) error {
	crs, ok := m.Courses[courseID]
	if !ok {
		return apperr.ErrCourseNotFound
	}
	updatedCrs, err := updateFn(ctx, &crs)
	if err != nil {
		return err
	}
	m.Courses[updatedCrs.ID()] = *updatedCrs
	return nil
}

type AcademicsService struct {
	Teachers map[string]bool
	Students map[string]bool
	Groups   map[string]bool
}

func (m *AcademicsService) TeacherExists(teacherID string) error {
	if m.Teachers[teacherID] {
		return nil
	}
	return apperr.ErrTeacherDoesntExist
}

func (m *AcademicsService) StudentExists(studentID string) error {
	if m.Students[studentID] {
		return nil
	}
	return apperr.ErrStudentDoesntExist
}

func (m *AcademicsService) GroupExists(groupID string) error {
	if m.Groups[groupID] {
		return nil
	}
	return apperr.ErrGroupDoesntExist
}
