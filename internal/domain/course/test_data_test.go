package course_test

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/domain/course"
)

func TestNewTestData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		InputData   string
		OutputData  string
		ExpectedErr error
	}{
		{
			Name:       "valid_test_data_parameters",
			InputData:  "1 + 1",
			OutputData: "2",
		},
		{
			Name:        "input_data_too_long",
			InputData:   strings.Repeat("x", 1001),
			OutputData:  "xxx",
			ExpectedErr: course.ErrTestInputDataTooLong,
		},
		{
			Name:        "output_data_too_long",
			InputData:   "xxx",
			OutputData:  strings.Repeat("x", 1001),
			ExpectedErr: course.ErrTestOutputDataTooLong,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			testData, err := course.NewTestData(c.InputData, c.OutputData)

			if c.ExpectedErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.ExpectedErr))

				return
			}
			require.NoError(t, err)
			require.Equal(t, c.InputData, testData.InputData())
			require.Equal(t, c.OutputData, testData.OutputData())
		})
	}
}

func TestTestData_IsZero(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name         string
		TestData     course.TestData
		ShouldBeZero bool
	}{
		{
			Name:         "should_be_zero",
			TestData:     course.TestData{},
			ShouldBeZero: true,
		},
		{
			Name:         "should_not_be_zero",
			TestData:     course.MustNewTestData("2 * 2", "4"),
			ShouldBeZero: false,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ShouldBeZero, c.TestData.IsZero())
		})
	}
}
