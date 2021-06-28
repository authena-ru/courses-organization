package command

import (
	"context"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type AddStudentHandler struct {
	coursesRepository coursesRepository
	academicsService  academicsService
}

func NewAddStudentHandler(repository coursesRepository, service academicsService) AddStudentHandler {
	if repository == nil {
		panic("coursesRepository is nil")
	}
	if service == nil {
		panic("academicsService is nil")
	}
	return AddStudentHandler{
		coursesRepository: repository,
		academicsService:  service,
	}
}

func (h AddStudentHandler) Handle(ctx context.Context, cmd app.AddStudentCommand) error {
	err := h.coursesRepository.UpdateCourse(ctx, cmd.CourseID, h.addStudent(cmd))
	return errors.Wrapf(
		err,
		"adding student #%s to course #%s by academic #%s",
		cmd.StudentID, cmd.CourseID, cmd.Academic.ID(),
	)
}

func (h AddStudentHandler) addStudent(cmd app.AddStudentCommand) UpdateFunction {
	return func(ctx context.Context, crs *course.Course) (*course.Course, error) {
		if err := h.academicsService.StudentExists(ctx, cmd.StudentID); err != nil {
			return nil, err
		}
		if err := crs.AddStudents(cmd.Academic, cmd.StudentID); err != nil {
			return nil, err
		}
		return crs, nil
	}
}
