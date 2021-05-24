package course_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

func TestNewPeriod(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name              string
		AcademicStartYear int
		AcademicEndYear   int
		Semester          course.Semester
		ShouldBeErr       bool
	}{
		{
			Name:              "valid_course_period",
			AcademicStartYear: 2021,
			AcademicEndYear:   2022,
			Semester:          course.FirstSemester,
		},
		{
			Name:              "academic_start_year_after_end",
			AcademicStartYear: 2024,
			AcademicEndYear:   2023,
			Semester:          course.SecondSemester,
			ShouldBeErr:       true,
		},
		{
			Name:              "academic_year_duration_over_year",
			AcademicStartYear: 2022,
			AcademicEndYear:   2024,
			Semester:          course.FirstSemester,
			ShouldBeErr:       true,
		},
		{
			Name:              "academic_start_year_equals_end_year",
			AcademicStartYear: 2023,
			AcademicEndYear:   2023,
			Semester:          course.FirstSemester,
			ShouldBeErr:       true,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			period, err := course.NewPeriod(c.AcademicStartYear, c.AcademicEndYear, c.Semester)
			if c.ShouldBeErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, c.AcademicStartYear, period.AcademicStartYear())
			require.Equal(t, c.AcademicEndYear, period.AcademicEndYear())
			require.Equal(t, c.Semester, period.Semester())
		})
	}
}

func TestPeriod_IsZero(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name         string
		Period       course.Period
		ShouldBeZero bool
	}{
		{
			Name:         "should_not_be_zero",
			Period:       course.MustNewPeriod(2021, 2022, course.SecondSemester),
			ShouldBeZero: false,
		},
		{
			Name:         "should_be_zero",
			Period:       course.Period{},
			ShouldBeZero: true,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ShouldBeZero, c.Period.IsZero())
		})
	}
}
