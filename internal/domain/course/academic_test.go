package course_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/domain/course"
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
			AcademicType: course.TeacherType,
		},
		{
			Name:         "empty_academic_id",
			AcademicType: course.StudentType,
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
			Academic:     course.MustNewAcademic("academic-id", course.TeacherType),
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
