package course_test

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestNewTestPoint(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name                  string
		Description           string
		Variants              []string
		CorrectVariantNumbers []int
		ExpectedErr           error
	}{
		{
			Name:                  "valid_test_point_parameters",
			Description:           "a or b?",
			Variants:              []string{"a", "b"},
			CorrectVariantNumbers: []int{0},
		},
		{
			Name:                  "test_point_description_too_long",
			Description:           strings.Repeat("x", 501),
			Variants:              []string{"x1", "x2"},
			CorrectVariantNumbers: []int{1},
			ExpectedErr:           course.ErrTestPointDescriptionTooLong,
		},
		{
			Name:                  "empty_test_point_variants",
			Description:           "1 + 1 = 2 or 3?",
			Variants:              nil,
			CorrectVariantNumbers: []int{0},
			ExpectedErr:           course.ErrEmptyTestPointVariants,
		},
		{
			Name:                  "empty_test_point_correct_variants",
			Description:           "2 * 2 = 4 or 2 * 2 = 5?",
			Variants:              []string{"4", "5"},
			CorrectVariantNumbers: nil,
			ExpectedErr:           course.ErrEmptyTestPointCorrectVariants,
		},
		{
			Name:                  "too_much_test_point_correct_variants",
			Description:           "Unsigned integer type: uint or int?",
			Variants:              []string{"uint", "int"},
			CorrectVariantNumbers: []int{0, 1, 2},
			ExpectedErr:           course.ErrTooMuchTestPointCorrectVariants,
		},
		{
			Name:                  "test_point_correct_variant_number_less_than_0",
			Description:           "Golang is OOP language",
			Variants:              []string{"Yes", "No"},
			CorrectVariantNumbers: []int{-1, 0},
			ExpectedErr:           course.ErrInvalidTestPointVariantNumber,
		},
		{
			Name:                  "test_point_correct_variant_number_more_than_last_variant",
			Description:           "Spring is awesome",
			Variants:              []string{"Yes", "No"},
			CorrectVariantNumbers: []int{1, 2},
			ExpectedErr:           course.ErrInvalidTestPointVariantNumber,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			tp, err := course.NewTestPoint(c.Description, c.Variants, c.CorrectVariantNumbers)

			if c.ExpectedErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.ExpectedErr))
				return
			}
			require.NoError(t, err)
			require.Equal(t, c.Description, tp.Description())
			require.Equal(t, c.Variants, tp.Variants())
			require.Equal(t, c.CorrectVariantNumbers, tp.CorrectVariantNumbers())
		})
	}
}

func TestTestPoint_IsZero(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name         string
		TestPoint    course.TestPoint
		ShouldBeZero bool
	}{
		{
			Name:         "should_be_zero",
			TestPoint:    course.TestPoint{},
			ShouldBeZero: true,
		},
		{
			Name:         "should_not_be_zero",
			TestPoint:    course.MustNewTestPoint("Django is cool", []string{"Yes", "No"}, []int{1}),
			ShouldBeZero: false,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ShouldBeZero, c.TestPoint.IsZero())
		})
	}
}
