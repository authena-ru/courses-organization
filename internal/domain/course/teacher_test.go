package course_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestCourse_AddCollaborators(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name     string
		Academic course.Academic
		IsErr    func(err error) bool
	}{
		{
			Name:     "creator_can_add_collaborators",
			Academic: course.MustNewAcademic("creator-id", course.TeacherType),
		},
		{
			Name:     "collaborator_can_add_collaborators",
			Academic: course.MustNewAcademic("collaborator-id", course.TeacherType),
		},
		{
			Name:     "student_cant_add_collaborators",
			Academic: course.MustNewAcademic("student-id", course.StudentType),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "not_course_teacher_cant_add_collaborators",
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

			err := crs.AddCollaborators(c.Academic, "collaborator1-id", "collaborator2-id")
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(
				t,
				[]string{"collaborator-id", "collaborator1-id", "collaborator2-id"},
				crs.Collaborators(),
			)
		})
	}
}

func TestCourse_RemoveCollaborators(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name     string
		Academic course.Academic
		IsErr    func(err error) bool
	}{
		{
			Name:     "creator_can_remove_collaborators",
			Academic: course.MustNewAcademic("creator-id", course.TeacherType),
		},
		{
			Name:     "collaborator_cant_remove_collaborators",
			Academic: course.MustNewAcademic("collaborator-id", course.TeacherType),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "student_cant_remove_collaborators",
			Academic: course.MustNewAcademic("student-id", course.StudentType),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "not_course_teacher_cant_remove_collaborators",
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
				t,
				creator,
				WithStudents("student-id"),
				WithCollaborators("collaborator-id", "collaborator-to-remove-id"),
			)

			err := crs.RemoveCollaborator(c.Academic, "collaborator-to-remove-id")
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, []string{"collaborator-id"}, crs.Collaborators())
		})
	}
}
