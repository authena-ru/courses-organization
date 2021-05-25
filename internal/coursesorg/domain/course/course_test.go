package course_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/authena-ru/courses-organization/internal/coursesorg/domain/course"
)

func TestNewCourse(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name        string
		Params      course.CreationCourseParams
		ShouldBeErr bool
		PossibleErr error
	}{
		{
			Name: "valid_course_creation_params",
			Params: course.CreationCourseParams{
				ID:        "course-id",
				CreatorID: "teacher-id",
				Title:     "Awesome Go in backend",
				Period:    course.MustNewPeriod(2022, 2023, course.SecondSemester),
				Started:   true,
			},
			ShouldBeErr: false,
		},
		{
			Name: "empty_course_id",
			Params: course.CreationCourseParams{
				CreatorID: "teacher-id",
				Title:     "Programming architecture",
				Period:    course.MustNewPeriod(2021, 2022, course.FirstSemester),
				Started:   false,
			},
			ShouldBeErr: true,
			PossibleErr: course.ErrEmptyCourseID,
		},
		{
			Name: "empty_course_creator_id",
			Params: course.CreationCourseParams{
				ID:      "course-id",
				Title:   "JavaScript in browser",
				Period:  course.MustNewPeriod(2023, 2024, course.SecondSemester),
				Started: true,
			},
			ShouldBeErr: true,
			PossibleErr: course.ErrEmptyCreatorID,
		},
		{
			Name: "empty_course_title",
			Params: course.CreationCourseParams{
				ID:        "course-id",
				CreatorID: "creator-id",
				Period:    course.MustNewPeriod(2024, 2025, course.FirstSemester),
				Started:   false,
			},
			ShouldBeErr: true,
			PossibleErr: course.ErrEmptyCourseTitle,
		},
		{
			Name: "zero_course_period",
			Params: course.CreationCourseParams{
				ID:        "course-id",
				CreatorID: "creator-id",
				Title:     "Nice React, Awesome Angular",
				Started:   true,
			},
			ShouldBeErr: true,
			PossibleErr: course.ErrZeroCoursePeriod,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs, err := course.NewCourse(c.Params)

			if c.ShouldBeErr {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.PossibleErr))
				return
			}
			require.NoError(t, err)
			require.Equal(t, c.Params.ID, crs.ID())
			require.Equal(t, c.Params.CreatorID, crs.CreatorID())
			require.Equal(t, c.Params.Title, crs.Title())
			require.Equal(t, c.Params.Period, crs.Period())
			require.Equal(t, c.Params.Started, crs.Started())
		})
	}
}
