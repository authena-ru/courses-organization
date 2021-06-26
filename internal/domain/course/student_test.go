package course_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestCourse_AddStudents(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name     string
		Academic course.Academic
		IsErr    func(err error) bool
	}{
		{
			Name:     "creator_can_add_students",
			Academic: course.MustNewAcademic("creator-id", course.TeacherType),
		},
		{
			Name:     "collaborator_can_add_students",
			Academic: course.MustNewAcademic("collaborator-id", course.TeacherType),
		},
		{
			Name:     "student_cant_add_students",
			Academic: course.MustNewAcademic("student-id", course.StudentType),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "not_course_teacher_cant_add_students",
			Academic: course.MustNewAcademic("another-teacher-id", course.TeacherType),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			creator := course.MustNewAcademic("creator-id", course.TeacherType)
			crs := NewCourse(t, creator, WithStudents("student-id"), WithCollaborators("collaborator-id"))
			err := crs.AddStudents(c.Academic, "student1-id", "student2-id")
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, []string{"student-id", "student1-id", "student2-id"}, crs.Students())
		})
	}
}

func TestCourse_RemoveStudents(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name     string
		Academic course.Academic
		IsErr    func(err error) bool
	}{
		{
			Name:     "creator_can_remove_students",
			Academic: course.MustNewAcademic("creator-id", course.TeacherType),
		},
		{
			Name:     "collaborator_can_remove_students",
			Academic: course.MustNewAcademic("collaborator-id", course.TeacherType),
		},
		{
			Name:     "student_cant_remove_students",
			Academic: course.MustNewAcademic("student-id", course.StudentType),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "not_course_teacher_cant_remove_students",
			Academic: course.MustNewAcademic("another-teacher-id", course.TeacherType),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			creator := course.MustNewAcademic("creator-id", course.TeacherType)
			crs := NewCourse(
				t, creator,
				WithStudents("student-id", "student-to-remove-id"),
				WithCollaborators("collaborator-id"),
			)
			err := crs.RemoveStudent(c.Academic, "student-to-remove-id")
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, []string{"student-id"}, crs.Students())
		})
	}
}
