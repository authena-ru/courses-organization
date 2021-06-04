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
		Params      course.CreationParams
		ExpectedErr error
	}{
		{
			Name: "valid_course_creation_params",
			Params: course.CreationParams{
				ID:            "course-id",
				Creator:       course.MustNewAcademic("creator-id", course.Teacher),
				Title:         "Awesome Go in backend",
				Period:        course.MustNewPeriod(2022, 2023, course.SecondSemester),
				Started:       true,
				Collaborators: []string{"collaborator-1-id", "collaborator-2-id"},
				Students:      []string{"student-id"},
			},
		},
		{
			Name: "empty_course_id",
			Params: course.CreationParams{
				Creator: course.MustNewAcademic("creator-id", course.Teacher),
				Title:   "Programming architecture",
				Period:  course.MustNewPeriod(2021, 2022, course.FirstSemester),
				Started: false,
			},
			ExpectedErr: course.ErrEmptyCourseID,
		},
		{
			Name: "zero_course_creator",
			Params: course.CreationParams{
				ID:      "course-id",
				Title:   "JavaScript in browser",
				Period:  course.MustNewPeriod(2023, 2024, course.SecondSemester),
				Started: true,
			},
			ExpectedErr: course.ErrZeroCreator,
		},
		{
			Name: "student_cant_create_course",
			Params: course.CreationParams{
				ID:      "course-id",
				Creator: course.MustNewAcademic("student-id", course.Student),
				Title:   "Assembly",
				Period:  course.MustNewPeriod(2020, 2021, course.FirstSemester),
				Started: false,
			},
			ExpectedErr: course.ErrNotTeacherCantCreateCourse,
		},
		{
			Name: "empty_course_title",
			Params: course.CreationParams{
				ID:      "course-id",
				Creator: course.MustNewAcademic("creator-id", course.Teacher),
				Period:  course.MustNewPeriod(2024, 2025, course.FirstSemester),
				Started: false,
			},
			ExpectedErr: course.ErrEmptyCourseTitle,
		},
		{
			Name: "zero_course_period",
			Params: course.CreationParams{
				ID:      "course-id",
				Creator: course.MustNewAcademic("creator-id", course.Teacher),
				Title:   "Nice React, Awesome Angular",
				Started: true,
			},
			ExpectedErr: course.ErrZeroCoursePeriod,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			crs, err := course.NewCourse(c.Params)

			if c.ExpectedErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, c.ExpectedErr))
				return
			}
			require.NoError(t, err)
			require.Equal(t, c.Params.ID, crs.ID())
			require.Equal(t, c.Params.Creator.ID(), crs.CreatorID())
			require.Equal(t, c.Params.Title, crs.Title())
			require.Equal(t, c.Params.Period, crs.Period())
			require.Equal(t, c.Params.Started, crs.Started())
			require.ElementsMatch(t, c.Params.Students, crs.Students())
			require.ElementsMatch(t, c.Params.Collaborators, crs.Collaborators())
		})
	}
}
