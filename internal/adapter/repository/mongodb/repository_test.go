package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/authena-ru/courses-organization/internal/adapter/repository/mongodb"
	"github.com/authena-ru/courses-organization/internal/domain/course"
)

type CoursesRepositoryTestSuite struct {
	suite.Suite
	MongoTestFixtures
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
	repo := s.newCoursesRepository()

	testCases := []struct {
		Name           string
		CoursesFactory func() *course.Course
	}{
		{
			Name: "course_without_collaborators and students",
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
	}

	for i := range testCases {
		c := testCases[i]

		s.Run(c.Name, func() {
			expectedCourse := c.CoursesFactory()

			ctx := context.Background()

			err := repo.AddCourse(ctx, expectedCourse)
			s.Require().NoError(err)

			s.requirePersistedCourseEquals(expectedCourse, repo)
		})
	}

	s.removeAllCourses(repo)
}

func (s *CoursesRepositoryTestSuite) newCoursesRepository() *mongodb.CoursesRepository {
	return mongodb.NewCoursesRepository(s.db)
}

func (s *CoursesRepositoryTestSuite) removeAllCourses(repo *mongodb.CoursesRepository) {
	err := repo.RemoveAllCourses(context.Background())
	s.Require().NoError(err)
}

func (s *CoursesRepositoryTestSuite) requirePersistedCourseEquals(
	expectedCourse *course.Course,
	repo *mongodb.CoursesRepository,
) {
	persistedCourse, err := repo.GetCourse(context.Background(), expectedCourse.ID())
	s.Require().NoError(err)

	s.Require().Equal(expectedCourse, persistedCourse)
}
