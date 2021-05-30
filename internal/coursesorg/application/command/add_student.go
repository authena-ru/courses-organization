package command

import (
	"context"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

type AddStudentCommand struct {
	Teacher   course.Academic
	CourseID  string
	StudentID string
}

type AddStudentHandler struct {
	coursesRepository coursesRepository
	academicsService  academicsService
}

func NewAddStudentHandler(repository coursesRepository, service academicsService) AddStudentHandler {
	if repository == nil {
		panic("CoursesRepository is nil")
	}
	if service == nil {
		panic("academicsService is nil")
	}
	return AddStudentHandler{
		coursesRepository: repository,
		academicsService:  service,
	}
}

// Handle is AddStudentCommand handler.
// Adds one student to course, returns error.
func (h AddStudentHandler) Handle(ctx context.Context, cmd AddStudentCommand) error {
	studentToAdd, err := course.NewAcademic(cmd.StudentID, course.Student)
	if err != nil {
		return err
	}
	if err := h.academicsService.AcademicExists(studentToAdd); err != nil {
		return err
	}
	return h.coursesRepository.UpdateCourse(
		ctx,
		cmd.CourseID,
		cmd.Teacher,
		func(_ context.Context, crs *course.Course) (*course.Course, error) {
			if err := crs.AddStudents(cmd.Teacher, studentToAdd.ID()); err != nil {
				return nil, err
			}
			return crs, nil
		})
}
