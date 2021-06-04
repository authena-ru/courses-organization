package course_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

func TestCourse_AddCollaborators(t *testing.T) {
	t.Parallel()
	var (
		creatorID            = "creator-id"
		collaboratorID       = "collaborator-id"
		studentID            = "student-id"
		collaboratorIDsToAdd = []string{"collaborator-1-id", "collaborator-2-id"}
	)
	crs := course.MustNewCourse(course.CreationParams{
		ID:            "creator-id",
		Creator:       course.MustNewAcademic(creatorID, course.Teacher),
		Title:         "ASP.NET in C#",
		Period:        course.MustNewPeriod(2021, 2022, course.SecondSemester),
		Collaborators: []string{collaboratorID},
		Students:      []string{studentID},
	})
	testCases := []struct {
		Name     string
		Course   course.Course
		Academic course.Academic
		IsErr    func(err error) bool
	}{
		{
			Name:     "creator_can_add_collaborators",
			Course:   *crs,
			Academic: course.MustNewAcademic(creatorID, course.Teacher),
		},
		{
			Name:     "collaborator_can_add_collaborators",
			Course:   *crs,
			Academic: course.MustNewAcademic(collaboratorID, course.Teacher),
		},
		{
			Name:     "student_cant_add_collaborators",
			Course:   *crs,
			Academic: course.MustNewAcademic(studentID, course.Student),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "not_course_teacher_cant_add_collaborators",
			Course:   *crs,
			Academic: course.MustNewAcademic("another-teacher-id", course.Teacher),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := c.Course.AddCollaborators(c.Academic, collaboratorIDsToAdd...)
			if err != nil {
				require.True(t, c.IsErr(err))
				return
			}
			totalCollaborators := append(collaboratorIDsToAdd, collaboratorID)
			require.Len(t, crs.Collaborators(), len(totalCollaborators))
			require.ElementsMatch(t, totalCollaborators, crs.Collaborators())
		})
	}
}

func TestCourse_RemoveCollaborators(t *testing.T) {
	t.Parallel()
	var (
		creatorID               = "creator-id"
		collaboratorID          = "collaborator-id"
		studentID               = "student-id"
		collaboratorIDsToRemove = []string{"collaborator-1-id", "collaborator-2-id"}
	)
	crs := course.MustNewCourse(course.CreationParams{
		ID:            "course-id",
		Creator:       course.MustNewAcademic(creatorID, course.Teacher),
		Title:         "GraphQL",
		Period:        course.MustNewPeriod(2025, 2026, course.FirstSemester),
		Collaborators: append(collaboratorIDsToRemove, collaboratorID),
		Students:      []string{studentID},
	})
	testCases := []struct {
		Name     string
		Course   course.Course
		Academic course.Academic
		IsErr    func(err error) bool
	}{
		{
			Name:     "creator_can_remove_collaborators",
			Course:   *crs,
			Academic: course.MustNewAcademic(creatorID, course.Teacher),
		},
		{
			Name:     "collaborator_cant_remove_collaborators",
			Course:   *crs,
			Academic: course.MustNewAcademic(collaboratorID, course.Teacher),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "student_cant_remove_collaborators",
			Course:   *crs,
			Academic: course.MustNewAcademic(studentID, course.Student),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "not_course_teacher_cant_remove_collaborators",
			Course:   *crs,
			Academic: course.MustNewAcademic("another-teacher-id", course.Teacher),
			IsErr:    course.IsAcademicCantEditCourseError,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := c.Course.RemoveCollaborators(c.Academic, collaboratorIDsToRemove...)
			if err != nil {
				require.True(t, c.IsErr(err))
				return
			}
			totalCollaborators := []string{collaboratorID}
			require.Len(t, crs.Collaborators(), len(totalCollaborators))
			require.ElementsMatch(t, totalCollaborators, crs.Collaborators())
		})
	}
}
