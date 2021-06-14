package course_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestCourse_AddStudents(t *testing.T) {
	t.Parallel()
	var (
		creatorID       = "creator-id"
		collaboratorID  = "collaborator-id"
		studentID       = "student-id"
		studentIDsToAdd = []string{"student-1-id", "student-2-id"}
	)
	params := course.CreationParams{
		ID:            "course-id",
		Creator:       course.MustNewAcademic(creatorID, course.Teacher),
		Title:         "SQL databases",
		Period:        course.MustNewPeriod(2021, 2022, course.FirstSemester),
		Started:       true,
		Collaborators: []string{collaboratorID},
		Students:      []string{studentID},
	}
	testCases := []struct {
		Name     string
		Academic course.Academic
		IsErr    func(err error) bool
	}{
		{
			Name:     "creator_can_add_students",
			Academic: course.MustNewAcademic(creatorID, course.Teacher),
		},
		{
			Name:     "collaborator_can_add_students",
			Academic: course.MustNewAcademic(collaboratorID, course.Teacher),
		},
		{
			Name:     "student_cant_add_students",
			Academic: course.MustNewAcademic(studentID, course.Student),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "not_course_teacher_cant_add_students",
			Academic: course.MustNewAcademic("another-teacher-id", course.Teacher),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(params)
			err := crs.AddStudents(c.Academic, studentIDsToAdd...)
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			totalStudents := append(studentIDsToAdd, studentID)
			require.ElementsMatch(t, totalStudents, crs.Students())
		})
	}
}

func TestCourse_RemoveStudents(t *testing.T) {
	t.Parallel()
	var (
		creatorID         = "creator-id"
		collaboratorID    = "collaborator-id"
		studentID         = "student-id"
		studentIDToRemove = "student-to-remove-id"
	)
	params := course.CreationParams{
		ID:            "course-id",
		Creator:       course.MustNewAcademic(creatorID, course.Teacher),
		Title:         "TypeScript from JavaScript",
		Period:        course.MustNewPeriod(2023, 2024, course.FirstSemester),
		Collaborators: []string{collaboratorID},
		Students:      []string{studentID, studentIDToRemove},
	}
	testCases := []struct {
		Name     string
		Academic course.Academic
		IsErr    func(err error) bool
	}{
		{
			Name:     "creator_can_remove_students",
			Academic: course.MustNewAcademic(creatorID, course.Teacher),
		},
		{
			Name:     "collaborator_can_remove_students",
			Academic: course.MustNewAcademic(collaboratorID, course.Teacher),
		},
		{
			Name:     "student_cant_remove_students",
			Academic: course.MustNewAcademic(studentID, course.Student),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "not_course_teacher_cant_remove_students",
			Academic: course.MustNewAcademic("another-teacher-id", course.Teacher),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(params)
			err := crs.RemoveStudent(c.Academic, studentIDToRemove)
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			totalStudents := []string{studentID}
			require.ElementsMatch(t, totalStudents, crs.Students())
		})
	}
}
