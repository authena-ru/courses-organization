package mongodb_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/authena-ru/courses-organization/internal/app"

	"github.com/stretchr/testify/suite"

	"github.com/authena-ru/courses-organization/internal/adapter/repository/mongodb"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type CoursesRepositoryTestSuite struct {
	suite.Suite
	MongoTestFixtures

	repository *mongodb.CoursesRepository
}

func (s *CoursesRepositoryTestSuite) SetupTest() {
	s.repository = s.newCoursesRepository()
}

func (s *CoursesRepositoryTestSuite) TearDownTest() {
	err := s.repository.RemoveAllCourses(context.Background())
	s.Require().NoError(err)
}

func TestCoursesRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration tests are skipped")
	}

	suite.Run(t, &CoursesRepositoryTestSuite{
		MongoTestFixtures: MongoTestFixtures{t: t},
	})
}

func (s *CoursesRepositoryTestSuite) TestCoursesRepository_AddCourse() {
	testCases := []struct {
		Name           string
		CoursesFactory func() *course.Course
	}{
		{
			Name: "course_without_collaborators_and_students",
			CoursesFactory: func() *course.Course {
				return course.MustNewCourse(course.CreationParams{
					ID:      "8f0d3461-175e-4e34-aebe-f292d83c3957",
					Creator: course.MustNewAcademic("aec5b656-1028-472c-a2be-f4bd2b0424b1", course.TeacherType),
					Title:   "New brand course",
					Period:  course.MustNewPeriod(2023, 2024, course.FirstSemester),
					Started: false,
				})
			},
		},
		{
			Name: "course_wit_collaborators",
			CoursesFactory: func() *course.Course {
				return course.MustNewCourse(course.CreationParams{
					ID:      "1e2339f9-7d31-42c9-bbfc-302c1471f48b",
					Creator: course.MustNewAcademic("93fedf0c-ae0d-4301-81b5-1f2291bdfa32", course.TeacherType),
					Title:   "New brand course with collaborators",
					Period:  course.MustNewPeriod(2024, 2025, course.SecondSemester),
					Started: true,
					Collaborators: []string{
						"8452d33b-a739-42d5-8090-627e080176b0",
						"1793f030-fedf-48b1-a996-7f8a2ff37a6f",
						"86094876-1d58-4a10-a7c7-4c0cb4e6584d",
					},
				})
			},
		},
		{
			Name: "course_with_students",
			CoursesFactory: func() *course.Course {
				return course.MustNewCourse(course.CreationParams{
					ID:      "7a7d3715-a54d-40f7-a568-95fa823b597e",
					Creator: course.MustNewAcademic("32a1be5d-62d5-4554-8784-507e95679966", course.TeacherType),
					Title:   "New brand course with students",
					Period:  course.MustNewPeriod(2025, 2026, course.FirstSemester),
					Started: true,
					Students: []string{
						"c5206f35-32d7-4ae8-903f-c4a2cad4c325",
						"57219c45-f7dc-4ac1-8695-d8aa2bdeb151",
						"e4841ac4-2bb6-4dd4-aabe-0e2a735ee613",
						"e81347e9-55f1-44a2-8adc-e8d65dd8b5ae",
					},
				})
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		s.Run(c.Name, func() {
			expectedCourse := c.CoursesFactory()

			ctx := context.Background()

			err := s.repository.AddCourse(ctx, expectedCourse)
			s.Require().NoError(err)

			s.requirePersistedCourseEquals(expectedCourse)
		})
	}
}

func (s *CoursesRepositoryTestSuite) TestCoursesRepository_GetCourse() {
	existingCourse := course.MustNewCourse(course.CreationParams{
		ID:      "dfe52190-ed85-4130-aba7-409c0bc21a49",
		Creator: course.MustNewAcademic("7b07ce84-d8aa-454b-9d01-305b211fe730", course.TeacherType),
		Title:   "New brand existing course",
		Period:  course.MustNewPeriod(2027, 2028, course.SecondSemester),
		Started: false,
	})

	err := s.repository.AddCourse(context.Background(), existingCourse)
	s.Require().NoError(err)

	testCases := []struct {
		Name          string
		CourseIDToGet string
		ExpectedErr   error
	}{
		{
			Name:          "existing_course",
			CourseIDToGet: existingCourse.ID(),
		},
		{
			Name:          "non_existing_course",
			CourseIDToGet: "23c645b3-9c94-4d5b-86e5-c8a3f9a47db3",
			ExpectedErr:   app.ErrCourseDoesntExist,
		},
	}

	for i := range testCases {
		c := testCases[i]

		s.Run(c.Name, func() {
			givenCourse, err := s.repository.GetCourse(context.Background(), c.CourseIDToGet)

			if c.ExpectedErr != nil {
				s.Require().True(errors.Is(err, c.ExpectedErr))

				return
			}

			s.Require().NoError(err)
			s.Require().Equal(existingCourse, givenCourse)
		})
	}
}

func (s *CoursesRepositoryTestSuite) newCoursesRepository() *mongodb.CoursesRepository {
	return mongodb.NewCoursesRepository(s.db)
}

func (s *CoursesRepositoryTestSuite) requirePersistedCourseEquals(expectedCourse *course.Course) {
	persistedCourse, err := s.repository.GetCourse(context.Background(), expectedCourse.ID())
	s.Require().NoError(err)

	s.Require().Equal(expectedCourse, persistedCourse)
}
