package course_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestCourse_AddCollaborators(t *testing.T) {
	t.Parallel()
	var (
		creatorID            = "creator-id"
		collaboratorID       = "collaborator-id"
		studentID            = "student-id"
		collaboratorIDsToAdd = []string{"collaborator-1-id", "collaborator-2-id"}
	)
	params := course.CreationParams{
		ID:            "course-id",
		Creator:       course.MustNewAcademic(creatorID, course.Teacher),
		Title:         "ASP.NET in C#",
		Period:        course.MustNewPeriod(2021, 2022, course.SecondSemester),
		Collaborators: []string{collaboratorID},
		Students:      []string{studentID},
	}
	testCases := []struct {
		Name     string
		Academic course.Academic
		IsErr    func(err error) bool
	}{
		{
			Name:     "creator_can_add_collaborators",
			Academic: course.MustNewAcademic(creatorID, course.Teacher),
		},
		{
			Name:     "collaborator_can_add_collaborators",
			Academic: course.MustNewAcademic(collaboratorID, course.Teacher),
		},
		{
			Name:     "student_cant_add_collaborators",
			Academic: course.MustNewAcademic(studentID, course.Student),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "not_course_teacher_cant_add_collaborators",
			Academic: course.MustNewAcademic("another-teacher-id", course.Teacher),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(params)

			err := crs.AddCollaborators(c.Academic, collaboratorIDsToAdd...)
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			totalCollaborators := append(collaboratorIDsToAdd, collaboratorID)
			require.ElementsMatch(t, totalCollaborators, crs.Collaborators())
		})
	}
}

func TestCourse_RemoveCollaborators(t *testing.T) {
	t.Parallel()
	var (
		creatorID              = "creator-id"
		collaboratorID         = "collaborator-id"
		studentID              = "student-id"
		collaboratorIDToRemove = "collaborator-to-remove-id"
	)
	params := course.CreationParams{
		ID:            "course-id",
		Creator:       course.MustNewAcademic(creatorID, course.Teacher),
		Title:         "GraphQL",
		Period:        course.MustNewPeriod(2025, 2026, course.FirstSemester),
		Collaborators: []string{collaboratorID, collaboratorIDToRemove},
		Students:      []string{studentID},
	}
	testCases := []struct {
		Name     string
		Academic course.Academic
		IsErr    func(err error) bool
	}{
		{
			Name:     "creator_can_remove_collaborators",
			Academic: course.MustNewAcademic(creatorID, course.Teacher),
		},
		{
			Name:     "collaborator_cant_remove_collaborators",
			Academic: course.MustNewAcademic(collaboratorID, course.Teacher),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "student_cant_remove_collaborators",
			Academic: course.MustNewAcademic(studentID, course.Student),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "not_course_teacher_cant_remove_collaborators",
			Academic: course.MustNewAcademic("another-teacher-id", course.Teacher),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs := course.MustNewCourse(params)

			err := crs.RemoveCollaborator(c.Academic, collaboratorIDToRemove)
			if c.IsErr != nil {
				require.Error(t, err)
				require.True(t, c.IsErr(err))
				return
			}
			require.NoError(t, err)
			totalCollaborators := []string{collaboratorID}
			require.ElementsMatch(t, totalCollaborators, crs.Collaborators())
		})
	}
}
