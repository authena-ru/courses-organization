package course_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

func TestNewAcademic(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name         string
		AcademicID   string
		AcademicType course.AcademicType
		ExpectedErr  error
	}{
		{
			Name:         "valid_academic_creation_params",
			AcademicID:   "teacher-id",
			AcademicType: course.Teacher,
		},
		{
			Name:         "empty_academic_id",
			AcademicType: course.Student,
			ExpectedErr:  course.ErrEmptyAcademicID,
		},
		{
			Name:        "invalid_academic_type",
			AcademicID:  "some-id",
			ExpectedErr: course.ErrInvalidAcademicType,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			academic, err := course.NewAcademic(c.AcademicID, c.AcademicType)

			if c.ExpectedErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.ExpectedErr))
				return
			}

			require.NoError(t, err)
			require.Equal(t, c.AcademicID, academic.ID())
			require.Equal(t, c.AcademicType, academic.Type())
		})
	}
}

func TestAcademic_IsZero(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name         string
		Academic     course.Academic
		ShouldBeZero bool
	}{
		{
			Name:         "should_not_be_zero",
			Academic:     course.MustNewAcademic("academic-id", course.Teacher),
			ShouldBeZero: false,
		},
		{
			Name:         "should_be_zero",
			Academic:     course.Academic{},
			ShouldBeZero: true,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ShouldBeZero, c.Academic.IsZero())
		})
	}
}

func TestCourse_CanAcademicSee(t *testing.T) {
	t.Parallel()
	var (
		creatorID      = "creator-id"
		collaboratorID = "collaborator-id"
		studentID      = "student-id"
		isErrFreeFunc  = func(err error) bool {
			return err == nil
		}
	)
	crs := course.MustNewCourse(course.CreationCourseParams{
		ID:            "course-id",
		Creator:       course.MustNewAcademic(creatorID, course.Teacher),
		Title:         "Nice Python",
		Period:        course.MustNewPeriod(2020, 2021, course.SecondSemester),
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
			Name:     "creator_can_see_course",
			Course:   *crs,
			Academic: course.MustNewAcademic(creatorID, course.Teacher),
			IsErr:    isErrFreeFunc,
		},
		{
			Name:     "collaborator_can_see_course",
			Course:   *crs,
			Academic: course.MustNewAcademic(collaboratorID, course.Teacher),
			IsErr:    isErrFreeFunc,
		},
		{
			Name:     "student_can_see_course",
			Course:   *crs,
			Academic: course.MustNewAcademic(studentID, course.Student),
			IsErr:    isErrFreeFunc,
		},
		{
			Name:     "not_participant_cant_see_course",
			Course:   *crs,
			Academic: course.MustNewAcademic("not-participant-id", course.Teacher),
			IsErr:    course.IsAcademicCantSeeCourseError,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := c.Course.CanAcademicSee(c.Academic)
			require.True(t, c.IsErr(err))
		})
	}
}

func TestCourse_CanAcademicEditWithAccess(t *testing.T) {
	t.Parallel()
	var (
		creatorID      = "creator-id"
		collaboratorID = "collaborator-id"
		studentID      = "student-id"
		isErrFreeFunc  = func(err error) bool {
			return err == nil
		}
	)
	crs := course.MustNewCourse(course.CreationCourseParams{
		ID:            "course-id",
		Creator:       course.MustNewAcademic(creatorID, course.Teacher),
		Title:         "Architecture of programming system",
		Period:        course.MustNewPeriod(2021, 2022, course.FirstSemester),
		Collaborators: []string{collaboratorID},
		Students:      []string{studentID},
	})
	testsCases := []struct {
		Name     string
		Course   course.Course
		Academic course.Academic
		Access   course.Access
		IsErr    func(err error) bool
	}{
		{
			Name:     "creator_can_edit_course_with_creator_access",
			Course:   *crs,
			Academic: course.MustNewAcademic(creatorID, course.Teacher),
			Access:   course.CreatorAccess,
			IsErr:    isErrFreeFunc,
		},
		{
			Name:     "creator_can_edit_course_with_teacher_access",
			Course:   *crs,
			Academic: course.MustNewAcademic(creatorID, course.Teacher),
			Access:   course.TeacherAccess,
			IsErr:    isErrFreeFunc,
		},
		{
			Name:     "collaborator_can_edit_course_with_teacher_access",
			Course:   *crs,
			Academic: course.MustNewAcademic(collaboratorID, course.Teacher),
			Access:   course.TeacherAccess,
			IsErr:    isErrFreeFunc,
		},
		{
			Name:     "collaborator_cant_edit_course_with_creator_access",
			Course:   *crs,
			Academic: course.MustNewAcademic(collaboratorID, course.Teacher),
			Access:   course.CreatorAccess,
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "student_cant_edit_course_with_teacher_access",
			Course:   *crs,
			Academic: course.MustNewAcademic(studentID, course.Student),
			Access:   course.TeacherAccess,
			IsErr:    course.IsAcademicCantEditCourseError,
		},
		{
			Name:     "student_cant_edit_course_with_creator_access",
			Course:   *crs,
			Academic: course.MustNewAcademic(studentID, course.Student),
			Access:   course.CreatorAccess,
			IsErr:    course.IsAcademicCantEditCourseError,
		},
	}

	for i := range testsCases {
		c := testsCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := c.Course.CanAcademicEditWithAccess(c.Academic, c.Access)
			require.True(t, c.IsErr(err))
		})
	}
}

func TestAcademic_CanCreateCourse(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name        string
		Academic    course.Academic
		ExpectedErr error
	}{
		{
			Name:     "teacher_can_create_course",
			Academic: course.MustNewAcademic("teacher-id", course.Teacher),
		},
		{
			Name:        "student_cant_create_course",
			Academic:    course.MustNewAcademic("student-id", course.Student),
			ExpectedErr: course.ErrNotTeacherCantCreateCourse,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			err := c.Academic.CanCreateCourse()

			if c.ExpectedErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.ExpectedErr))
				return
			}
			require.NoError(t, err)
		})
	}
}
