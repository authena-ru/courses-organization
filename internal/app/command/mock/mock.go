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

type AcademicsService struct {
	teachers map[string]bool
	students map[string]bool
	groups   map[string]bool
}

func NewAcademicsService(teachers []string, students []string, groups []string) *AcademicsService {
	asm := &AcademicsService{
		teachers: make(map[string]bool, len(teachers)),
		students: make(map[string]bool, len(students)),
		groups:   make(map[string]bool, len(groups)),
	}
	for _, t := range teachers {
		asm.teachers[t] = true
	}

	for _, s := range students {
		asm.students[s] = true
	}

	for _, g := range groups {
		asm.groups[g] = true
	}

	return asm
}

func (m *AcademicsService) TeacherExists(_ context.Context, teacherID string) error {
	if m.teachers[teacherID] {
		return nil
	}

	return app.ErrTeacherDoesntExist
}

func (m *AcademicsService) StudentExists(_ context.Context, studentID string) error {
	if m.students[studentID] {
		return nil
	}

	return app.ErrStudentDoesntExist
}

func (m *AcademicsService) GroupExists(_ context.Context, groupID string) error {
	if m.groups[groupID] {
		return nil
	}

	return app.ErrGroupDoesntExist
}
