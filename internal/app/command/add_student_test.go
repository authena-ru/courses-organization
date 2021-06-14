package command_test

import (
	"context"
	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/app/command"
	"github.com/authena-ru/courses-organization/internal/app/command/mock"
	"github.com/authena-ru/courses-organization/internal/domain/course"
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
	addCourse := func(crs *course.Course) *mock.CoursesRepository {
		return mock.NewCoursesRepository(crs)
	}
	addStudent := func() *mock.AcademicsService {
		return mock.NewAcademicsService(nil, []string{studentID}, nil)
	}
	testCases := []struct {
		Name                     string
		Command                  command.AddStudentCommand
		PrepareCoursesRepository func(crs *course.Course) *mock.CoursesRepository
		PrepareAcademicsService  func() *mock.AcademicsService
		IsErr                    func(err error) bool
	}{
		{
			Name: "add_student",
			Command: command.AddStudentCommand{
				Academic:  creator,
				CourseID:  courseID,
				StudentID: studentID,
			},
			PrepareCoursesRepository: addCourse,
			PrepareAcademicsService:  addStudent,
		},
		{
			Name: "dont_add_when_teacher_cant_edit_course",
			Command: command.AddStudentCommand{
				Academic:  course.MustNewAcademic("other-teacher-id", course.Teacher),
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
				Academic:  creator,
				CourseID:  courseID,
				StudentID: studentID,
			},
			PrepareCoursesRepository: addCourse,
			PrepareAcademicsService: func() *mock.AcademicsService {
				return mock.NewAcademicsService(nil, nil, nil)
			},
			IsErr: func(err error) bool {
				return errors.Is(err, app.ErrStudentDoesntExist)
			},
		},
		{
			Name: "dont_add_when_course_doesnt_exist",
			Command: command.AddStudentCommand{
				Academic:  creator,
				CourseID:  courseID,
				StudentID: studentID,
			},
			PrepareCoursesRepository: func(_ *course.Course) *mock.CoursesRepository {
				return mock.NewCoursesRepository()
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
			coursesRepository := c.PrepareCoursesRepository(crs)
			academicsService := c.PrepareAcademicsService()
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
