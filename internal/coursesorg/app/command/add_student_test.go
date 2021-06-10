package command_test

import (
	"context"
	"github.com/authena-ru/courses-organization/internal/coursesorg/app"
	"github.com/authena-ru/courses-organization/internal/coursesorg/app/command"
	"github.com/authena-ru/courses-organization/internal/coursesorg/app/command/mock"
	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddStudentHandler_Handle(t *testing.T) {
	t.Parallel()
	var (
		courseID  = "course-id"
		studentID = "student-id"
		creator   = course.MustNewAcademic("creator-id", course.Teacher)
	)
	addCourse := func(crs *course.Course, crm *mock.CoursesRepository) {
		crm.Courses = map[string]course.Course{crs.ID(): *crs}
	}
	addStudent := func(asm *mock.AcademicsService) {
		asm.Students = map[string]bool{studentID: true}
	}
	testCases := []struct {
		Name                     string
		Command                  command.AddStudentCommand
		PrepareCoursesRepository func(crs *course.Course, crm *mock.CoursesRepository)
		PrepareAcademicsService  func(asm *mock.AcademicsService)
		IsErr                    func(err error) bool
	}{
		{
			Name: "add_student",
			Command: command.AddStudentCommand{
				Teacher:   creator,
				CourseID:  courseID,
				StudentID: studentID,
			},
			PrepareCoursesRepository: addCourse,
			PrepareAcademicsService:  addStudent,
		},
		{
			Name: "dont_add_when_teacher_cant_edit_course",
			Command: command.AddStudentCommand{
				Teacher:   course.MustNewAcademic("other-teacher-id", course.Teacher),
				CourseID:  courseID,
				StudentID: studentID,
			},
			PrepareCoursesRepository: addCourse,
			PrepareAcademicsService:  addStudent,
			IsErr:                    course.IsAcademicCantEditCourseError,
		},
		{
			Name: "dont_add_when_student_doesnt_exist",
			Command: command.AddStudentCommand{
				Teacher:   creator,
				CourseID:  courseID,
				StudentID: studentID,
			},
			PrepareCoursesRepository: addCourse,
			PrepareAcademicsService: func(asm *mock.AcademicsService) {
				asm.Students = make(map[string]bool)
			},
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrStudentDoesntExist)
			},
		},
		{
			Name: "dont_add_when_course_doesnt_exist",
			Command: command.AddStudentCommand{
				Teacher:   creator,
				CourseID:  courseID,
				StudentID: studentID,
			},
			PrepareCoursesRepository: func(_ *course.Course, crm *mock.CoursesRepository) {
				crm.Courses = make(map[string]course.Course)
			},
			PrepareAcademicsService: addStudent,
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrCourseDoesntExist)
			},
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(course.CreationParams{
				ID:      courseID,
				Creator: creator,
				Title:   "Math",
				Period:  course.MustNewPeriod(2028, 2029, course.SecondSemester),
			})
			coursesRepository := &mock.CoursesRepository{}
			c.PrepareCoursesRepository(crs, coursesRepository)
			academicsService := &mock.AcademicsService{}
			c.PrepareAcademicsService(academicsService)
			handler := command.NewAddStudentHandler(coursesRepository, academicsService)

			err := handler.Handle(context.Background(), c.Command)

			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				require.NotContains(t, crs.Students(), c.Command.StudentID)
				return
			}
			require.NoError(t, err)
			require.Contains(t, crs.Students(), c.Command.StudentID)
		})
	}
}
