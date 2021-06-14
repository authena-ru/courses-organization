package course_test

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestNewDeadline(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name               string
		ExcellentGradeTime time.Time
		GoodGradeTime      time.Time
		ExpectedErr        error
	}{
		{
			Name:               "valid_time_for_creation_deadline",
			ExcellentGradeTime: time.Date(2023, time.September, 22, 0, 0, 0, 0, time.Local),
			GoodGradeTime:      time.Date(2023, time.October, 02, 0, 0, 0, 0, time.Local),
			ExpectedErr:        nil,
		},
		{
			Name:               "zero_excellent_grade_time",
			ExcellentGradeTime: time.Time{},
			GoodGradeTime:      time.Date(2021, time.November, 01, 0, 0, 0, 0, time.Local),
			ExpectedErr:        course.ErrZeroExcellentGradeTime,
		},
		{
			Name:               "zero_good_grade_time",
			ExcellentGradeTime: time.Date(2022, time.April, 02, 0, 0, 0, 0, time.Local),
			GoodGradeTime:      time.Time{},
			ExpectedErr:        course.ErrZeroGoodGradeTime,
		},
		{
			Name:               "excellent_grade_time_after_good",
			ExcellentGradeTime: time.Date(2022, time.September, 28, 0, 0, 0, 0, time.Local),
			GoodGradeTime:      time.Date(2022, time.September, 01, 0, 0, 0, 0, time.Local),
			ExpectedErr:        course.ErrExcellentGradeTimeAfterGood,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			deadline, err := course.NewDeadline(c.ExcellentGradeTime, c.GoodGradeTime)

			if c.ExpectedErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.ExpectedErr))
				return
			}

			require.NoError(t, err)
			require.Equal(t, c.ExcellentGradeTime, deadline.ExcellentGradeTime())
			require.Equal(t, c.GoodGradeTime, deadline.GoodGradeTime())
		})
	}
}

func TestDeadline_IsZero(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name         string
		Deadline     course.Deadline
		ShouldBeZero bool
	}{
		{
			Name:         "should_be_zero",
			Deadline:     course.Deadline{},
			ShouldBeZero: true,
		},
		{
			Name: "should_not_be_zero",
			Deadline: course.MustNewDeadline(
				time.Date(2021, time.October, 02, 0, 0, 0, 0, time.Local),
				time.Date(2021, time.October, 12, 0, 0, 0, 0, time.Local),
			),
			ShouldBeZero: false,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ShouldBeZero, c.Deadline.IsZero())
		})
	}
}
