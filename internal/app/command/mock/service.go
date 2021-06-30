package mock

import (
	"context"

	"github.com/authena-ru/courses-organization/internal/app"
)

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
